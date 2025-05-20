// Package cooldown is used to make cooldowns in go. It has basic and nonbasic implementations. My
// goal was to make it well-documented.
//
// Basic is very easy designed and have no features. Unlike others:
// ValuedHandler and CoolDown, they're adding self to global Processor (Proc var)
//
// In theory, user should bbe allowed to do custom processors, but there are unexported interfaces
// and methods, that'll be exported in future updates.
//
// Processor is doing most work, because it checks, if cooldown is expired, and if it is, stops it
// immediately.
package cooldown
