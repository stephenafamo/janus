# Janus: The many faced framework

Janus is a modular framework for Go web applications.

It is made up of several interfaces (hence "many faced") which can be picked-and-used as required. Each interface will come with several implementations, so the user does not have to write their own implementation unless necessary.

A work in progress...

## Philosophy

The general advice in the Go community is to use the standard library as much as possible and not jump into frameworks.

However, most Go frameworks I've seen are heavily opinionated which means that unless you being with them, or convert almost all your code to be compatible with it, you can't really use the framework. It also makes it hard to remove if neccessary.

The purpose of **Janus** is to have a set of interfaces that rely mostly on the `stdlib`, with various implementations of these interfaces.

These interfaces should help with:

* Routing
* Rendering
* Authentication
* Migration
* Object storage

...more in the future

A user can then pick what parts they want to use and drop it into their application.

The dependence on the `stdlib` also makes it super easy to remove any unnecessary part of the framework and either use a compatible library, or the `stdlib` itself.