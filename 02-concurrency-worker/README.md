# 高并发任务系统 (goroutine + worker pool)

本项目是Go实战入门路线的第二阶段，目标是掌握Go的核心优势——并发模型。

## 项目目标

- Worker Pool（可动态扩容）
- goroutine 任务调度
- channel 队列
- 任务状态跟踪
- HTTP + WebSocket 推送任务进度
- 防止 goroutine 泄漏
- context 控制超时与取消

## 学习重点

- Go 与 Java 线程模型的根本差异
- CSP 并发模型
- 如何写不 panic、不泄漏、不阻塞的 worker 系统

## 启动方式

待完善...