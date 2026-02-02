// Package cooldown is used to make cooldowns in go. It has basic and nonbasic
// implementations.
//
// Basic is very easy designed and have no features. Unlike others:
// ValuedHandler and CoolDown: they're starting timer via time.AfterFunc, that
// cancels itself, when cooldown expires.
package cooldown
