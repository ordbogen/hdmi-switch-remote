var gulp = require('gulp'),
    concat = require('gulp-concat'),
    filter = require('gulp-filter'),
    size = require('gulp-size'),
    jade = require('gulp-jade'),
    stylus = require('gulp-stylus'),
    mainBowerFiles = require('main-bower-files');

gulp.task("default", () => {
  gulp
    .src("app/app.js")
    .pipe(gulp.dest("public"));

  gulp
    .src("app/index.jade")
    .pipe(jade())
    .pipe(gulp.dest("public"));

  gulp
    .src("app/hello.styl")
    .pipe(stylus())
    .pipe(gulp.dest("public"));

  gulp
    .src(mainBowerFiles())
    .pipe(filter("*.js"))
    .pipe(size({
        showFiles: true
    }))
    .pipe(concat("vendor.js"))
    .pipe(gulp.dest("public"));

  gulp
    .src("bower_components/angular-material/angular-material.min.css")
    .pipe(gulp.dest("public"));

  gulp
    .src("bower_components/components-font-awesome/css/font-awesome.min.css")
    .pipe(gulp.dest("public"));

  gulp
    .src("app/hdmi-switch-icon-192.png")
    .pipe(gulp.dest("public"));

  gulp
    .src("bower_components/components-font-awesome/fonts/*")
    .pipe(gulp.dest("public/fonts"));
});
