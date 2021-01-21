package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseConfig(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test Parseconfig", t, func() {
		config := ParseConfig([]byte(`[main]
version   = "1.0.0"
tag       = false
[[operate]]
location = "pyproject.toml"
search   = "version = \"{}\""
replace  = "version = \"{}\""`))

		Convey("version = 1.0.0", func() {
			So(config.VersionHelper.Version, ShouldEqual, "1.0.0")
		})
		Convey("tag flag = false", func() {
			So(config.VersionHelper.TagFlag, ShouldEqual, false)
		})
		Convey("2.1.0 -> 3.0.0 compatible", func() {
			Convey("serialize = {version}-{banner}", func() {
				So(config.VersionHelper.Serialize, ShouldEqual, "{version}-{banner}")
			})
		})

	})
}

func TestIsPureVersion(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test IsPureVersion", t, func() {

		Convey("Normal: 1.0.0 = true", func() {
			So(IsPureVersion("1.0.0"), ShouldEqual, true)
		})
		Convey("Invalid: invalid1.0.0 = false", func() {
			So(IsPureVersion("evil1.0.0"), ShouldEqual, false)
		})
	})
}

func TestParseVersionBySerialize(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test ParseVersionBySerialize", t, func() {

		Convey("Normal: {version}-{banner} 1.0.0-alpha", func() {
			resultMap, err := parseVersionBySerialize("{version}-{banner}", "1.0.0-alpha")
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("version = 1.0.0", func() {
				So(resultMap["version"], ShouldEqual, "1.0.0")
			})

			Convey("banner = alpha", func() {
				So(resultMap["banner"], ShouldEqual, "alpha")
			})
		})
		Convey("Without banner: {version}-{banner} 1.0.0", func() {
			resultMap, err := parseVersionBySerialize("{version}-{banner}", "1.0.0")
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("version = 1.0.0", func() {
				So(resultMap["version"], ShouldEqual, "1.0.0")
			})

			Convey("banner = ", func() {
				So(resultMap["banner"], ShouldBeEmpty)
			})
		})
		Convey("Extra serialize: [{version}]-{banner} [1.0.0]-alpha", func() {
			resultMap, err := parseVersionBySerialize("[{version}]-{banner}", "[1.0.0]-alpha")
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("version = 1.0.0", func() {
				So(resultMap["version"], ShouldEqual, "1.0.0")
			})

			Convey("banner = alpha", func() {
				So(resultMap["banner"], ShouldEqual, "alpha")
			})
		})
		Convey("Extra serialize without banner: [{version}]-{banner} [1.0.0]", func() {
			resultMap, err := parseVersionBySerialize("[{version}]-{banner}", "[1.0.0]")
			Convey("err = Parse Failed", func() {
				So(err.Error(), ShouldEqual, "Parse failed")
			})
			Convey("version = ", func() {
				So(resultMap["version"], ShouldEqual, "")
			})

			Convey("banner = ", func() {
				So(resultMap["banner"], ShouldEqual, "")
			})
		})
		Convey("serialize Without {banner}: {version} 1.0.0-alpha", func() {
			resultMap, err := parseVersionBySerialize("{version}", "1.0.0-alpha")
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("version = 1.0.0-alpha", func() {
				So(resultMap["version"], ShouldEqual, "1.0.0-alpha")
			})

			Convey("banner = ", func() {
				So(resultMap["banner"], ShouldBeEmpty)
			})
		})
		Convey("Invalid serialize: {} 1.0.0-alpha", func() {
			resultMap, err := parseVersionBySerialize("{}", "1.0.0-alpha")
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("version = ", func() {
				So(resultMap["version"], ShouldBeEmpty)
			})

			Convey("banner = ", func() {
				So(resultMap["banner"], ShouldBeEmpty)
			})
		})
	})
}

func TestGenerateVersion(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test GenerateVersion", t, func() {

		Convey("Without banner: 1.0.0, ,{version}-{banner}", func() {
			version, err := GenerateVersion("1.0.0", "", "{version}-{banner}")
			Convey("version = 1.0.0", func() {
				So(version, ShouldEqual, "1.0.0")
			})
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})

		})
		Convey("Normal: 1.0.0,alpha,{version}-{banner}", func() {
			version, err := GenerateVersion("1.0.0", "alpha", "{version}-{banner}")
			Convey("version = 1.0.0-alpha", func() {
				So(version, ShouldEqual, "1.0.0-alpha")
			})
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("Extra serialize: 1.0.0, ,[{version}]-{banner}", func() {
			version, err := GenerateVersion("1.0.0", "", "[{version}]-{banner}")
			Convey("version = 1.0.0", func() {
				So(version, ShouldEqual, "1.0.0")
			})
			Convey("err = nil", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("Invalid version: invalid,alpha,{version}-{banner}", func() {
			version, err := GenerateVersion("invalid", "alpha", "{version}-{banner}")
			Convey("version = ", func() {
				So(version, ShouldEqual, "")
			})
			Convey("err = Serialize invalid", func() {
				So(err.Error(), ShouldEqual, "Version invalid")
			})
		})
		Convey("Invalid serialize: 1.0.0,alpha,{invalid}-{banner}", func() {
			version, err := GenerateVersion("1.0.0", "alpha", "{invalid}-{banner}")
			Convey("version = ", func() {
				So(version, ShouldEqual, "")
			})
			Convey("err = Serialize invalid", func() {
				So(err.Error(), ShouldEqual, "Serialize invalid")
			})
		})
	})
}
