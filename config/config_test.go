package config

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/go-test/deep"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	a := assert.New(t)
	c := InitConfig("proj", "reg", "buck", "table", "prof", "me@foo.example", "0.99.0")
	a.Equal("prof", *c.Defaults.Common.Backend.Profile)
	a.Equal("prof", *c.Defaults.Providers.AWS.Profile)
	a.Equal("reg", *c.Defaults.Providers.AWS.Region)
	a.Equal("reg", *c.Defaults.Providers.AWS.Region)
	a.Equal("0.99.0", *c.Defaults.Providers.AWS.Version)
	a.Equal("buck", *c.Defaults.Common.Backend.Bucket)
	a.Equal("table", *c.Defaults.Common.Backend.DynamoTable)
	a.Equal("me@foo.example", *c.Defaults.Owner)
	a.Equal("proj", *c.Defaults.Project)
	a.Equal(defaultTerraformVersion.String(), *c.Defaults.TerraformVersion)
}

func Test_detectVersion(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"implicit 1", args{[]byte(`{}`)}, 1, false},
		{"explicit 1", args{[]byte(`{"version": 1}`)}, 1, false},
		{"explicit 2", args{[]byte(`{"version": 2}`)}, 2, false},
		{"err", args{[]byte(`{`)}, 0, true},
	}
	a := assert.New(t)
	fs, _, err := util.TestFs()
	a.NoError(err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int
			var err error
			switch tt.want {
			case 1:
				afero.WriteFile(fs, "fogg.json", tt.args.b, 0644)
				got, err = detectVersion(tt.args.b, fs, "fogg.json")
			case 2:
				afero.WriteFile(fs, "fogg.yml", tt.args.b, 0644)
				got, err = detectVersion(tt.args.b, fs, "fogg.yml")
			default:
				got, err = detectVersion(tt.args.b, fs, "fogg.yml")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("detectVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("detectVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func intptr(i int) *int {
	return &i
}

func jsonNumberPtr(i int) *json.Number {
	j := json.Number(strconv.Itoa(i))
	return &j
}

func strptr(s string) *string {
	return &s
}

func boolptr(b bool) *bool {
	return &b
}

func TestUpgradeConfigVersion(t *testing.T) {
	a := assert.New(t)
	b, e := util.TestFile("v1_full")
	a.NoError(e)
	v1Full, e := v1.ReadConfig(b)
	a.NoError(e)
	v2Full := &v2.Config{
		Version: 2,
		Defaults: v2.Defaults{
			Common: v2.Common{
				Backend: &v2.Backend{
					Bucket:      util.StrPtr("the-bucket"),
					DynamoTable: util.StrPtr("the-table"),
					Profile:     util.StrPtr("czi"),
					Region:      util.StrPtr("us-west-2"),
				},

				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID:         jsonNumberPtr(1),
						Profile:           strptr("czi"),
						Region:            strptr("us-west-1"),
						Version:           strptr("0.1.0"),
						AdditionalRegions: []string{"us-east-1"},
					},
				},
				Owner:            util.StrPtr("default@example.com"),
				Project:          util.StrPtr("test-project"),
				ExtraVars:        map[string]string{"foo": "bar"},
				TerraformVersion: util.StrPtr("0.11.0"),
				Tools: &v2.Tools{
					TfLint: &v1.TfLint{
						Enabled: boolptr(true),
					},
					TravisCI: &v1.TravisCI{
						CommonCI: v1.CommonCI{
							Enabled:        boolptr(true),
							AWSIAMRoleName: util.StrPtr("travis-role"),
							TestBuckets:    intptr(13),
							Buildevents:    boolptr(true),
						},
					},
				},
			},
		},

		Accounts: map[string]v2.Account{
			"foo": v2.Account{
				Common: v2.Common{
					Backend: &v2.Backend{
						Bucket:      util.StrPtr("foo-bucket"),
						DynamoTable: util.StrPtr("foo-table"),
						Profile:     util.StrPtr("czi-foo"),
						Region:      util.StrPtr("us-west-foo1"),
					},

					Providers: &v2.Providers{
						AWS: &v2.AWSProvider{
							AccountID:         jsonNumberPtr(2),
							Profile:           strptr("czi-foo"),
							Region:            strptr("us-west-foo1"),
							Version:           strptr("0.12.0"),
							AdditionalRegions: []string{"us-east-foo2"},
						},
					},
					Owner:            util.StrPtr("foo@example.com"),
					Project:          util.StrPtr("foo-project"),
					ExtraVars:        map[string]string{"foo": "foo"},
					TerraformVersion: util.StrPtr("0.12.0"),
				},
			},
			"bar": v2.Account{
				Common: v2.Common{
					Backend: &v2.Backend{
						Bucket:      util.StrPtr("bar-bucket"),
						DynamoTable: util.StrPtr("bar-table"),
						Profile:     util.StrPtr("czi-bar"),
						Region:      util.StrPtr("us-west-bar1"),
					},

					Providers: &v2.Providers{
						AWS: &v2.AWSProvider{
							AccountID:         jsonNumberPtr(3),
							Profile:           strptr("czi-bar"),
							Region:            strptr("us-west-bar1"),
							Version:           strptr("0.13.0"),
							AdditionalRegions: []string{"us-east-bar2"},
						},
					},
					Owner:            util.StrPtr("bar@example.com"),
					Project:          util.StrPtr("bar-project"),
					ExtraVars:        map[string]string{"foo": "bar"},
					TerraformVersion: util.StrPtr("0.13.0"),
				},
			},
		},

		Envs: map[string]v2.Env{
			"stage": v2.Env{
				Common: v2.Common{
					Backend: &v2.Backend{
						Bucket:      util.StrPtr("stage-bucket"),
						DynamoTable: util.StrPtr("stage-table"),
						Profile:     util.StrPtr("czi-stage"),
						Region:      util.StrPtr("us-west-stage1"),
					},

					Providers: &v2.Providers{
						AWS: &v2.AWSProvider{
							AccountID:         jsonNumberPtr(4),
							Profile:           strptr("czi-stage"),
							Region:            strptr("us-west-stage1"),
							Version:           strptr("0.14.0"),
							AdditionalRegions: []string{"us-east-stage2"},
						},
					},
					Owner:            util.StrPtr("stage@example.com"),
					Project:          util.StrPtr("stage-project"),
					TerraformVersion: util.StrPtr("0.14.0"),
					ExtraVars:        map[string]string{"foo": "stage"},
				},
				Components: map[string]v2.Component{
					"cloud-env": v2.Component{
						Common: v2.Common{
							Backend: &v2.Backend{
								Bucket:      util.StrPtr("env-bucket"),
								DynamoTable: util.StrPtr("env-table"),
								Profile:     util.StrPtr("czi-env"),
								Region:      util.StrPtr("us-west-env1"),
							},

							Providers: &v2.Providers{
								AWS: &v2.AWSProvider{
									AccountID:         jsonNumberPtr(5),
									Profile:           strptr("czi-env"),
									Region:            strptr("us-west-env1"),
									Version:           strptr("0.15.0"),
									AdditionalRegions: []string{"us-east-env2"},
								},
							},
							Owner:            util.StrPtr("env@example.com"),
							Project:          util.StrPtr("env-project"),
							TerraformVersion: util.StrPtr("0.15.0"),
							ExtraVars:        map[string]string{"foo": "env"},
						},
						ModuleSource: strptr("github.com/chanzuckerberg/fogg-test-module"),
					},
					"helm": {},
				},
			},
		},

		Modules: map[string]v1.Module{
			"module1": v1.Module{TerraformVersion: strptr("0.15.0")},
		},
		Plugins: v1.Plugins{
			CustomPlugins: map[string]*plugins.CustomPlugin{
				"plugin": &plugins.CustomPlugin{
					URL:    "https://example.com/plugin.tgz",
					Format: plugins.TypePluginFormatTar,
					TarConfig: &plugins.TarConfig{
						StripComponents: 7,
					},
				},
			},
			TerraformProviders: map[string]*plugins.CustomPlugin{
				"provider": &plugins.CustomPlugin{
					URL:    "https://example.com/provider.tgz",
					Format: plugins.TypePluginFormatTar,
					TarConfig: &plugins.TarConfig{
						StripComponents: 7,
					},
				},
			},
		},
	}

	type args struct {
		c1 *v1.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *v2.Config
		wantErr bool
	}{
		{"v1 full", args{v1Full}, v2Full, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpgradeConfigVersion(tt.args.c1)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpgradeConfigVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := deep.Equal(tt.want, got); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestFindAndReadConfig(t *testing.T) {
	a := assert.New(t)

	fs := func(m map[string][]byte) (afero.Fs, error) {
		fs, _, err := util.TestFs()
		a.NoError(err)

		for k, v := range m {
			f, e := fs.OpenFile(k, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if e != nil {
				return nil, e
			}
			_, e = f.Write(v)
			if e != nil {
				return nil, e
			}
		}
		return fs, nil
	}

	v1, e := util.TestFile("v1_full")
	a.NoError(e)

	v2, e := util.TestFile("v2_minimal_valid")

	a.NoError(e)

	f1, e := fs(map[string][]byte{
		"config.json": v1,
	})
	a.NoError(e)
	defer f1.RemoveAll(".") //nolint

	f2, e := fs(map[string][]byte{
		"config.json": v2,
	})
	a.NoError(e)
	defer f2.RemoveAll(".") //nolint

	fErr, e := fs(map[string][]byte{
		"config.json": []byte(`{"version": 7}`),
	})
	a.NoError(e)
	defer fErr.RemoveAll(".") //nolint

	_, e = FindAndReadConfig(f1, "config.json")
	a.NoError(e)

	_, e = FindAndReadConfig(f2, "config.json")
	a.NoError(e)

	_, e = FindAndReadConfig(fErr, "config.json")
	a.Error(e)

}
