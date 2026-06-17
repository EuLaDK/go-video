# AGENTS.md

**Tradeoff:** These guidelines prioritize caution over speed. Use judgment for trivial tasks.

## 1. Think Before Coding

**State assumptions, surface tradeoffs, and ask when uncertain.**

- Explicitly state assumptions. Ask if instructions are unclear or if multiple interpretations exist.
- Propose simpler approaches and raise questions regarding overcomplication before implementing.

## 2. Simplicity First

**Write the minimum code required to solve the problem.**

- Avoid speculative features, unrequested flexibility, or single-use abstractions.
- Keep code concise; rewrite if it can be significantly shorter.

## 3. Surgical Changes

**Only modify what is strictly necessary.**

- Match existing style and avoid unrelated refactoring, formatting, or cleaning of adjacent code.
- **Document Functions:** Add clear Chinese comments explaining the purpose and parameters of any new functions you write.
- Mention unrelated dead code rather than deleting it.
- Remove only the unused imports, variables, or functions created by your changes.

## 4. Goal-Driven Execution

**Define clear success criteria and verify changes systematically.**

- Transform tasks into verifiable goals (e.g., write reproducing tests before fixing bugs).
- For multi-step tasks, outline a brief step-by-step verification plan:

```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

## 5. Project-Specific Boundaries

**Verify locally and reuse existing patterns.**

- Limit verification and type-checking to the modified folder or direct scope. Do not run full builds unless requested.
- Prefer Naive UI components over custom controls.

---

**Success criteria:** Minimal diffs, less overcomplication, and clarifying questions asked before writing code.
