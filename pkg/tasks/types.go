package tasks

type Task func() error

type TaskDecorator = func(task Task) Task
