## Understand graceful shutdown
   ┌──────────────┐
   │  main()      │
   └────┬─────────┘
        │
        ├─ initializeDependencies()
        ├─ signal listener goroutine ───┐
        ├─ RunWorker() goroutine        │ (listens for ctx cancel)
        ├─ startServer() goroutine      │ (listens for ctx cancel)
        └─ wg.Wait() <──────────────────┘ waits for both
