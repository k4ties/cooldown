// Package cooldown is used to make cooldowns in go. It has basic and nonbasic implementations. My
// goal was to make it very customizable and well-documented.
//
// Basic is very easy designed and have no features. Unlike others:
// ValuedHandler and CoolDown, they're adding self to global Processor (Proc var)
//
// In theory, user should be allowed to do custom processors without any issues, so, you can try,
// if you want :)
// Custom processors is usable in cases, where you don't want to spawn new goroutine, and you
// already have a ticker, where you can call processor routines.
//
// Processor is doing most work, because it checks, if cooldown is expired, and if it is, stops it
// immediately with StopCauseExpired stop cause.
// If Close method on Processor will be called, while some Processable is active, stop cause will
// become StopCauseClosed.
package cooldown
