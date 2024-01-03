package finder

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/manager/fakes"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/util"
)

func TestFinder_Find(t *testing.T) {
	testData, err := filepath.Abs("../../testdata")
	if err != nil {
		require.NoError(t, err)
	}

	testCases := []struct {
		name            string
		pluginDirs      []string
		pluginClass     plugins.Class
		expectedBundles []*plugins.FoundBundle
		err             error
	}{
		{
			name:       "Dir with single plugin",
			pluginDirs: []string{filepath.Join(testData, "valid-v2-signature")},
			expectedBundles: []*plugins.FoundBundle{
				{
					Primary: plugins.FoundPlugin{
						JSONData: plugins.JSONData{
							ID:   "test-datasource",
							Type: plugins.TypeDataSource,
							Name: "Test",
							Info: plugins.Info{
								Author: plugins.InfoLink{
									Name: "Will Browne",
									URL:  "https://willbrowne.com",
								},
								Description: "Test",
								Version:     "1.0.0",
							},
							Dependencies: plugins.Dependencies{
								GrafanaVersion: "*",
								Plugins:        []plugins.Dependency{},
							},
							State:      plugins.ReleaseStateAlpha,
							Backend:    true,
							Executable: "test",
						},
						FS: mustNewStaticFSForTests(t, filepath.Join(testData, "valid-v2-signature/plugin")),
					},
				},
			},
		},
		{
			name:       "Dir with nested plugins",
			pluginDirs: []string{"../../testdata/duplicate-plugins"},
			expectedBundles: []*plugins.FoundBundle{
				{
					Primary: plugins.FoundPlugin{
						JSONData: plugins.JSONData{
							ID:   "test-app",
							Type: plugins.TypeDataSource,
							Name: "Parent",
							Info: plugins.Info{
								Author: plugins.InfoLink{
									Name: "Grafana Labs",
									URL:  "http://grafana.com",
								},
								Description: "Parent plugin",
								Version:     "1.0.0",
								Updated:     "2020-10-20",
							},
							Dependencies: plugins.Dependencies{
								GrafanaVersion: "*",
								Plugins:        []plugins.Dependency{},
							},
						},
						FS: mustNewStaticFSForTests(t, filepath.Join(testData, "duplicate-plugins/nested")),
					},
					Children: []*plugins.FoundPlugin{
						{
							JSONData: plugins.JSONData{
								ID:   "test-app",
								Type: plugins.TypeDataSource,
								Name: "Child",
								Info: plugins.Info{
									Author: plugins.InfoLink{
										Name: "Grafana Labs",
										URL:  "http://grafana.com",
									},
									Description: "Child plugin",
									Version:     "1.0.0",
									Updated:     "2020-10-20",
								},
								Dependencies: plugins.Dependencies{
									GrafanaVersion: "*",
									Plugins:        []plugins.Dependency{},
								},
							},
							FS: mustNewStaticFSForTests(t, filepath.Join(testData, "duplicate-plugins/nested/nested")),
						},
					},
				},
			},
		},
		{
			name:       "Dir with single plugin which has symbolic link root directory",
			pluginDirs: []string{"../../testdata/symbolic-plugin-dirs"},
			expectedBundles: []*plugins.FoundBundle{
				{
					Primary: plugins.FoundPlugin{
						JSONData: plugins.JSONData{
							ID:   "test-app",
							Type: plugins.TypeApp,
							Name: "Test App",
							Info: plugins.Info{
								Author: plugins.InfoLink{
									Name: "Test Inc.",
									URL:  "http://test.com",
								},
								Description: "Official Grafana Test App & Dashboard bundle",
								Version:     "1.0.0",
								Links: []plugins.InfoLink{
									{Name: "Project site", URL: "http://project.com"},
									{Name: "License & Terms", URL: "http://license.com"},
								},
								Updated: "2015-02-10",
								Logos: plugins.Logos{
									Small: "img/logo_small.png",
									Large: "img/logo_large.png",
								},
								Screenshots: []plugins.Screenshots{
									{Name: "img1", Path: "img/screenshot1.png"},
									{Name: "img2", Path: "img/screenshot2.png"},
								},
							},
							Dependencies: plugins.Dependencies{
								GrafanaVersion: "3.x.x",
								Plugins: []plugins.Dependency{
									{ID: "graphite", Type: "datasource", Name: "Graphite", Version: "1.0.0"},
									{ID: "graph", Type: "panel", Name: "Graph", Version: "1.0.0"},
								},
							},
							Includes: []*plugins.Includes{
								{
									Name: "Nginx Connections",
									Path: "dashboards/connections.json",
									Type: "dashboard",
									Role: "Viewer",
								},
								{
									Name: "Nginx Memory",
									Path: "dashboards/memory.json",
									Type: "dashboard",
									Role: "Viewer",
								},
								{Name: "Nginx Panel", Type: "panel", Role: "Viewer"},
								{Name: "Nginx Datasource", Type: "datasource", Role: "Viewer"},
							},
						},
						FS: mustNewStaticFSForTests(t, filepath.Join(testData, "includes-symlinks")),
					},
				},
			},
		},
		{
			name:       "Multiple plugin dirs",
			pluginDirs: []string{"../../testdata/duplicate-plugins", "../../testdata/invalid-v1-signature"},
			expectedBundles: []*plugins.FoundBundle{{
				Primary: plugins.FoundPlugin{
					JSONData: plugins.JSONData{
						ID:   "test-app",
						Type: plugins.TypeDataSource,
						Name: "Parent",
						Info: plugins.Info{
							Author: plugins.InfoLink{
								Name: "Grafana Labs",
								URL:  "http://grafana.com",
							},
							Description: "Parent plugin",
							Version:     "1.0.0",
							Updated:     "2020-10-20",
						},
						Dependencies: plugins.Dependencies{
							GrafanaVersion: "*",
							Plugins:        []plugins.Dependency{},
						},
					},
					FS: mustNewStaticFSForTests(t, filepath.Join(testData, "duplicate-plugins/nested")),
				},
				Children: []*plugins.FoundPlugin{
					{
						JSONData: plugins.JSONData{
							ID:   "test-app",
							Type: plugins.TypeDataSource,
							Name: "Child",
							Info: plugins.Info{
								Author: plugins.InfoLink{
									Name: "Grafana Labs",
									URL:  "http://grafana.com",
								},
								Description: "Child plugin",
								Version:     "1.0.0",
								Updated:     "2020-10-20",
							},
							Dependencies: plugins.Dependencies{
								GrafanaVersion: "*",
								Plugins:        []plugins.Dependency{},
							},
						},
						FS: mustNewStaticFSForTests(t, filepath.Join(testData, "duplicate-plugins/nested/nested")),
					},
				},
			},
				{
					Primary: plugins.FoundPlugin{
						JSONData: plugins.JSONData{
							ID:   "test-datasource",
							Type: plugins.TypeDataSource,
							Name: "Test",
							Info: plugins.Info{
								Author: plugins.InfoLink{
									Name: "Grafana Labs",
									URL:  "https://grafana.com",
								},
								Description: "Test",
							},
							Dependencies: plugins.Dependencies{
								GrafanaVersion: "*",
								Plugins:        []plugins.Dependency{},
							},
							State:   plugins.ReleaseStateAlpha,
							Backend: true,
						},
						FS: mustNewStaticFSForTests(t, filepath.Join(testData, "invalid-v1-signature/plugin")),
					},
				},
			},
		},
		{
			name:        "Plugin with dist folder (core class)",
			pluginDirs:  []string{filepath.Join(testData, "plugin-with-dist")},
			pluginClass: plugins.ClassCore,
			expectedBundles: []*plugins.FoundBundle{
				{
					Primary: plugins.FoundPlugin{
						JSONData: plugins.JSONData{
							ID:   "test-datasource",
							Type: plugins.TypeDataSource,
							Name: "Test",
							Info: plugins.Info{
								Author: plugins.InfoLink{
									Name: "Will Browne",
									URL:  "https://willbrowne.com",
								},
								Description: "Test",
								Version:     "1.0.0",
							},
							Dependencies: plugins.Dependencies{
								GrafanaVersion: "*",
								Plugins:        []plugins.Dependency{},
							},
							State:      plugins.ReleaseStateAlpha,
							Backend:    true,
							Executable: "test",
						},
						FS: mustNewStaticFSForTests(t, filepath.Join(testData, "plugin-with-dist/plugin/dist")),
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewLocalFinder(false, featuremgmt.WithFeatures(featuremgmt.FlagExternalCorePlugins))
			pluginBundles, err := f.Find(context.Background(), &fakes.FakePluginSource{
				PluginClassFunc: func(ctx context.Context) plugins.Class {
					return tc.pluginClass
				},
				PluginURIsFunc: func(ctx context.Context) []string {
					return tc.pluginDirs
				},
			})
			if (err != nil) && !errors.Is(err, tc.err) {
				t.Errorf("Find() error = %v, expected error %v", err, tc.err)
				return
			}

			// to ensure we can compare with expected
			sort.SliceStable(pluginBundles, func(i, j int) bool {
				return pluginBundles[i].Primary.JSONData.ID < pluginBundles[j].Primary.JSONData.ID
			})

			if !cmp.Equal(pluginBundles, tc.expectedBundles, fsComparer) {
				t.Fatalf("Result mismatch (-want +got):\n%s", cmp.Diff(pluginBundles, tc.expectedBundles, fsComparer))
			}
		})
	}
}

func TestFinder_getAbsPluginJSONPaths(t *testing.T) {
	t.Run("When scanning a folder that doesn't exists shouldn't return an error", func(t *testing.T) {
		origWalk := walk
		walk = func(path string, followSymlinks, detectSymlinkInfiniteLoop, followDistFolder bool, walkFn util.WalkFunc) error {
			return walkFn(path, nil, os.ErrNotExist)
		}
		t.Cleanup(func() {
			walk = origWalk
		})

		finder := NewLocalFinder(false, featuremgmt.WithFeatures())
		paths, err := finder.getAbsPluginJSONPaths("test", true)
		require.NoError(t, err)
		require.Empty(t, paths)
	})

	t.Run("When scanning a folder that lacks permission shouldn't return an error", func(t *testing.T) {
		origWalk := walk
		walk = func(path string, followSymlinks, detectSymlinkInfiniteLoop, followDistFolder bool, walkFn util.WalkFunc) error {
			return walkFn(path, nil, os.ErrPermission)
		}
		t.Cleanup(func() {
			walk = origWalk
		})

		finder := NewLocalFinder(false, featuremgmt.WithFeatures())
		paths, err := finder.getAbsPluginJSONPaths("test", true)
		require.NoError(t, err)
		require.Empty(t, paths)
	})

	t.Run("When scanning a folder that returns a non-handled error should return that error", func(t *testing.T) {
		origWalk := walk
		walk = func(path string, followSymlinks, detectSymlinkInfiniteLoop, followDistFolder bool, walkFn util.WalkFunc) error {
			return walkFn(path, nil, errors.New("random error"))
		}
		t.Cleanup(func() {
			walk = origWalk
		})

		finder := NewLocalFinder(false, featuremgmt.WithFeatures())
		paths, err := finder.getAbsPluginJSONPaths("test", true)
		require.Error(t, err)
		require.Empty(t, paths)
	})

	t.Run("should forward if the dist folder should be evaluated", func(t *testing.T) {
		origWalk := walk
		walk = func(path string, followSymlinks, detectSymlinkInfiniteLoop, followDistFolder bool, walkFn util.WalkFunc) error {
			if followDistFolder {
				return walkFn(path, nil, errors.New("unexpected followDistFolder"))
			}
			return walkFn(path, nil, filepath.SkipDir)
		}
		t.Cleanup(func() {
			walk = origWalk
		})

		finder := NewLocalFinder(false, featuremgmt.WithFeatures())
		paths, err := finder.getAbsPluginJSONPaths("test", false)
		require.ErrorIs(t, err, filepath.SkipDir)
		require.Empty(t, paths)
	})
}

func TestFinder_getAbsPluginJSONPaths_PluginClass(t *testing.T) {
	t.Run("When a dist folder exists as a direct child of the plugins path, it will always be resolved", func(t *testing.T) {
		dir, err := filepath.Abs("../../testdata/pluginRootWithDist")
		require.NoError(t, err)

		tcs := []struct {
			name       string
			followDist bool
			expected   []string
		}{
			{
				name:       "When followDistFolder is enabled, only the nested dist folder will be followed",
				followDist: true,
				expected: []string{
					filepath.Join(dir, "dist/plugin.json"),
				},
			},
			{
				name:       "When followDistFolder is disabled, no dist folders will be followed",
				followDist: false,
				expected: []string{
					filepath.Join(dir, "datasource/plugin.json"),
					filepath.Join(dir, "panel/src/plugin.json"),
				},
			},
		}
		for _, tc := range tcs {
			pluginBundles, err := NewLocalFinder(false, featuremgmt.WithFeatures()).getAbsPluginJSONPaths(dir, tc.followDist)
			require.NoError(t, err)

			sort.Strings(pluginBundles)
			require.Equal(t, tc.expected, pluginBundles)
		}
	})
}

var fsComparer = cmp.Comparer(func(fs1 plugins.FS, fs2 plugins.FS) bool {
	fs1Files, err := fs1.Files()
	if err != nil {
		panic(err)
	}
	fs2Files, err := fs2.Files()
	if err != nil {
		panic(err)
	}

	sort.SliceStable(fs1Files, func(i, j int) bool {
		return fs1Files[i] < fs1Files[j]
	})

	sort.SliceStable(fs2Files, func(i, j int) bool {
		return fs2Files[i] < fs2Files[j]
	})

	return cmp.Equal(fs1Files, fs2Files) && fs1.Base() == fs2.Base()
})

func mustNewStaticFSForTests(t *testing.T, dir string) plugins.FS {
	sfs, err := plugins.NewStaticFS(plugins.NewLocalFS(dir))
	require.NoError(t, err)
	return sfs
}
