```mermaid
flowchart TD
  A(HotKey Press Down)
  B[Listen For Audio]
  A-->B
  C[Threshold and VAT Filtering]
  B--Stream PCM Chunks-->C
  D[HotKey Release]
```
