package apply

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/container-registry/harbor-satellite/internal/groundcontrol/cli/common"
	"github.com/container-registry/harbor-satellite/internal/groundcontrol/cli/models"
	"github.com/container-registry/harbor-satellite/pkg/groundcontrol"
)

func applyGroupResource(rt *common.Runtime, resource models.YAMLResource) models.ApplyResult {
	groupSpec, ok := resource.Spec.(*models.GroupSpec)
	if !ok {
		return errorToResult(resource.Kind, resource.Metadata.Name, fmt.Errorf("invalid spec format"))
	}
	action := models.ActionNil

	resp, err := rt.Client().GetGroupWithResponse(context.Background(), resource.Metadata.Name)
	if err != nil {
		return errorToResult(resource.Kind, resource.Metadata.Name, fmt.Errorf("failed to get current state: %v", err))
	}
	if resp.StatusCode() == http.StatusNotFound {
		action = models.ActionCreated
	}
	if resp.StatusCode() != http.StatusOK {
		return errorToResult(resource.Kind, resource.Metadata.Name, fmt.Errorf("failed to get current state: %v", err))
	}

	if action == models.ActionNil {
		action = compareGroupResource(groupSpec, responseToSpec(resp.JSON200))
	}

	switch action {
	}

	return models.ApplyResult{} // TODO: Placeholder
}

func compareGroupResource(spec1, spec2 *models.GroupSpec) models.Action {
	changed := false

	if spec1.RegistryURL != spec2.RegistryURL {
		changed = true
	}
	if !groupsEqual(spec1.Projects, spec2.Projects) {
		changed = true
	}

	if changed {
		return models.ActionConfigured
	}

	return models.ActionUnchanged
}

func responseToSpec(grp *groundcontrol.GroupResponse) *models.GroupSpec {
	return &models.GroupSpec{
		RegistryURL: grp.RegistryURL,
		Projects:    grp.Projects,
	}
}

func groupsEqual(a, b []string) bool {
	a, b = slices.Clone(a), slices.Clone(b)
	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(a, b)
}
