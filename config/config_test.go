package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestConfig(t *testing.T)  {

	Convey("Given an environment with no environment variables", t, func() {
		os.Clearenv()
		Convey("When the config values are retrieved", func() {
			cfg, err := Get()
			Convey("The values should be set to the default values", func() {
				So(cfg.BindAddr, ShouldEqual, ":8080")
				So(cfg.MongoConfig.BindAddr, ShouldEqual, "localhost:27017")
				So(cfg.MongoConfig.Database, ShouldEqual, "bookStore")
				So(cfg.MongoConfig.Collection, ShouldEqual, "books")
			})
			Convey("And there should be no errors", func() {
				So(err, ShouldBeNil)
			})
		})
	} )
	
}
