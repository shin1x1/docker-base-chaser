package handler

import (
	"fmt"
	"os/exec"

	"github.com/Masterminds/semver/v3"
	"github.com/shin1x1/docker-base-chaser/internal/provider"
	"github.com/shin1x1/docker-base-chaser/internal/template"
)

type RootHandler struct {
	dryRun bool
}

func NewRootHandler(dryRun bool) *RootHandler {
	return &RootHandler{dryRun: dryRun}
}

func (r *RootHandler) Exec(targets *Targets) error {
	for _, t := range *targets {
		pd := provider.New(t.Provider, t.Image)
		for {
			if err := r.fetchLoop(pd, t); err != nil {
				return err
			}

			if t.Done() || !pd.HasNext() {
				break
			}
		}
	}

	return nil
}

func (r *RootHandler) fetchLoop(provider provider.Provider, target *Target) error {
	images, err := provider.Fetch()
	if err != nil {
		return err
	}

	if err := r.updateTags(images, target); err != nil {
		return err
	}

	for _, tag := range target.Tags {
		if !tag.CanExecute() {
			continue
		}

		// テンプレート置換
		if err := template.ExecFile(target.Template.Source, target.Template.Destination, tag.Tag); err != nil {
			return err
		}

		// コマンド実行
		cmdText, err := template.Exec(target.Command, tag.Tag, target.Template.Destination)
		if err != nil {
			return err
		}

		prefix := "execute"
		if r.dryRun {
			prefix = "dry-run"
		}
		fmt.Printf("\n[%s]: Tag=%s\n%s", prefix, tag.Tag, cmdText)

		if r.dryRun {
			return nil
		}

		cmd := exec.Command("sh", "-c", cmdText)
		if err := cmd.Run(); err != nil {
			return err
		}

		tag.Executed()
	}

	return nil
}

func (r *RootHandler) updateTags(images []*provider.Image, target *Target) error {
	for _, tag := range target.Tags {
		if !tag.CanMatch() {
			continue;
		}

		for _, img := range images {
			// pattern にマッチすること
			if !tag.MatchPattern(img.Tag) {
				continue
			}

			tag.Matched()

			// 更新日付
			if !tag.ShouldUpdate(img) {
				tag.NotExecuted()
				continue
			}

			// version にマッチすること
			v, err := semver.NewVersion(img.Tag)
			if err != nil {
				continue
			}
			v1, _ := v.SetMetadata("")
			v2, _ := v1.SetPrerelease("")
			if !tag.CheckVersion(&v2) {
				continue
			}

			tag.Update(img)
			break
		}
	}

	return nil
}
