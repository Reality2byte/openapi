package outputfilter

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/openapi/tools/cli/internal/openapi"
	"github.com/oasdiff/oasdiff/load"
	"github.com/stretchr/testify/require"
)

func TestOpenApiSpecMethods(t *testing.T) {
	loader := openapi.NewOpenAPI3().WithExcludedPrivatePaths()
	spec, err := loader.CreateOpenAPISpecFromPath("../../../test/data/changelog/manual_generated_test_spec.json")
	require.NoError(t, err)

	expectedConfig := map[string]*OperationConfig{
		"createCluster": {
			Path:                   "/api/atlas/v2/groups/{groupId}/clusters",
			HTTPMethod:             "POST",
			Tag:                    "Multi-Cloud Clusters",
			Sunset:                 "2025-06-01",
			ManualChangelogEntries: map[string]any{},
		},
		"deleteCluster": {
			Path:                   "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}",
			HTTPMethod:             "DELETE",
			Tag:                    "Multi-Cloud Clusters",
			Sunset:                 "2025-06-01",
			ManualChangelogEntries: map[string]any{},
		},
		"getCluster": {
			Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}",
			HTTPMethod: "GET",
			Tag:        "Multi-Cloud Clusters",
			Sunset:     "2025-06-01",
			ManualChangelogEntries: map[string]any{
				"2023-07-07": "July 7th changelog entry",
				"2023-08-08": "August announcement for /api/atlas/v2/clusters",
			},
		},
		"listClusters": {
			Path:                   "/api/atlas/v2/groups/{groupId}/clusters",
			HTTPMethod:             "GET",
			Tag:                    "Multi-Cloud Clusters",
			Sunset:                 "2025-06-01",
			ManualChangelogEntries: map[string]any{},
		},
		"updateCluster": {
			Path:                   "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}",
			HTTPMethod:             "PATCH",
			Tag:                    "Multi-Cloud Clusters",
			Sunset:                 "2025-06-01",
			ManualChangelogEntries: map[string]any{},
		},
	}

	operationConfigs := newOperationConfigFromSpec(spec)
	require.True(t, reflect.DeepEqual(expectedConfig, operationConfigs))
}

func TestOperationConfigs_Tag(t *testing.T) {
	tests := []struct {
		name     string
		configs  OperationConfigs
		expected string
	}{
		{
			name: "Tag from Revision",
			configs: OperationConfigs{
				Base:     &OperationConfig{Tag: "base-tag"},
				Revision: &OperationConfig{Tag: "revision-tag"},
			},
			expected: "revision-tag",
		},
		{
			name: "Tag from Base",
			configs: OperationConfigs{
				Base: &OperationConfig{Tag: "base-tag"},
			},
			expected: "base-tag",
		},
		{
			name: "No Tag",
			configs: OperationConfigs{
				Base:     nil,
				Revision: nil,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.configs.Tag(); got != tt.expected {
				t.Errorf("Tag() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOperationConfigs_Sunset(t *testing.T) {
	tests := []struct {
		name     string
		configs  OperationConfigs
		expected string
	}{
		{
			name: "Sunset from Revision",
			configs: OperationConfigs{
				Base:     &OperationConfig{Sunset: "2024-07-24"},
				Revision: &OperationConfig{Sunset: "2024-08-01"},
			},
			expected: "2024-08-01",
		},
		{
			name: "Sunset from Base",
			configs: OperationConfigs{
				Base: &OperationConfig{Sunset: "2024-07-24"},
			},
			expected: "2024-07-24",
		},
		{
			name: "No Sunset",
			configs: OperationConfigs{
				Base:     nil,
				Revision: nil,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.configs.Sunset(); got != tt.expected {
				t.Errorf("Sunset() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewEndpointsConfigGivenBaseAndRevision(t *testing.T) {
	tests := []struct {
		name         string
		baseSpec     *load.SpecInfo
		revisionSpec *load.SpecInfo
		expected     map[string]*OperationConfigs
	}{
		{
			name: "Merge base and revision with different operation IDs",
			baseSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/base-path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"base-tag"},
							},
						}),
					),
				},
			},
			revisionSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/revision-path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op2",
								Tags:        []string{"revision-tag"},
							},
						}),
					),
				},
			},
			expected: map[string]*OperationConfigs{
				"op1": {
					Base: &OperationConfig{
						Path:                   "/base-path",
						HTTPMethod:             "GET",
						Tag:                    "base-tag",
						Sunset:                 "",
						ManualChangelogEntries: map[string]any{},
					},
					Revision: nil,
				},
				"op2": {
					Base: nil,
					Revision: &OperationConfig{
						Path:                   "/revision-path",
						HTTPMethod:             "GET",
						Tag:                    "revision-tag",
						Sunset:                 "",
						ManualChangelogEntries: map[string]any{},
					},
				},
			},
		},
		{
			name: "Merge base and revision with the same operation IDs",
			baseSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"base-tag"},
								Extensions: map[string]any{
									"x-sunset": "2024-08-01",
								},
							},
						}),
					),
				},
			},
			revisionSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"revision-tag"},
								Extensions: map[string]any{
									"x-sunset": "2024-12-01",
								},
							},
						}),
					),
				},
			},
			expected: map[string]*OperationConfigs{
				"op1": {
					Base: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "base-tag",
						Sunset:                 "2024-08-01",
						ManualChangelogEntries: map[string]any{},
					},
					Revision: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "revision-tag",
						Sunset:                 "2024-12-01",
						ManualChangelogEntries: map[string]any{},
					},
				},
			},
		},
		{
			name: "Base with sunset extension only",
			baseSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"base-tag"},
								Extensions: map[string]any{
									"x-sunset": "2024-08-01",
								},
							},
						}),
					),
				},
			},
			revisionSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"revision-tag"},
							},
						}),
					),
				},
			},
			expected: map[string]*OperationConfigs{
				"op1": {
					Base: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "base-tag",
						Sunset:                 "2024-08-01",
						ManualChangelogEntries: map[string]any{},
					},
					Revision: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "revision-tag",
						Sunset:                 "", // Sunset not present in revision
						ManualChangelogEntries: map[string]any{},
					},
				},
			},
		},
		{
			name: "Revision with sunset extension only",
			baseSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"base-tag"},
							},
						}),
					),
				},
			},
			revisionSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"revision-tag"},
								Extensions: map[string]any{
									"x-sunset": "2024-12-01",
								},
							},
						}),
					),
				},
			},
			expected: map[string]*OperationConfigs{
				"op1": {
					Base: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "base-tag",
						ManualChangelogEntries: map[string]any{},
					},
					Revision: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "revision-tag",
						Sunset:                 "2024-12-01",
						ManualChangelogEntries: map[string]any{},
					},
				},
			},
		},
		{
			name: "Revision with x-xgen-changelog extension only",
			baseSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"base-tag"},
							},
						}),
					),
				},
			},
			revisionSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"revision-tag"},
								Extensions: map[string]any{
									"x-xgen-changelog": map[string]any{
										"2024-01-01": "Revision changelog entry.",
									},
								},
							},
						}),
					),
				},
			},
			expected: map[string]*OperationConfigs{
				"op1": {
					Base: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "base-tag",
						Sunset:                 "",
						ManualChangelogEntries: map[string]any{},
					},
					Revision: &OperationConfig{
						Path:       "/path",
						HTTPMethod: "GET",
						Tag:        "revision-tag",
						Sunset:     "",
						ManualChangelogEntries: map[string]any{
							"2024-01-01": "Revision changelog entry.",
						},
					},
				},
			},
		},
		{
			name: "Revision with multiple x-xgen-changelog extensions",
			baseSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"base-tag"},
							},
						}),
					),
				},
			},
			revisionSpec: &load.SpecInfo{
				Spec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath("/path", &openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "op1",
								Tags:        []string{"revision-tag"},
								Extensions: map[string]any{
									"x-xgen-changelog": map[string]any{
										"2024-01-01": "Revision changelog entry.",
										"2024-01-02": "Revision changelog entry.",
									},
								},
							},
						}),
					),
				},
			},
			expected: map[string]*OperationConfigs{
				"op1": {
					Base: &OperationConfig{
						Path:                   "/path",
						HTTPMethod:             "GET",
						Tag:                    "base-tag",
						Sunset:                 "",
						ManualChangelogEntries: map[string]any{},
					},
					Revision: &OperationConfig{
						Path:       "/path",
						HTTPMethod: "GET",
						Tag:        "revision-tag",
						Sunset:     "",
						ManualChangelogEntries: map[string]any{
							"2024-01-01": "Revision changelog entry.",
							"2024-01-02": "Revision changelog entry.",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewOperationConfigs(tt.baseSpec, tt.revisionSpec)
			if !reflect.DeepEqual(tt.expected, result) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
