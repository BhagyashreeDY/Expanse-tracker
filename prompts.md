# Project Build Logs: Interaction History

This document logs the core prompts and instructions used to build this engine. We kept this for transparency to show how the system evolved from a basic API to a production-ready settlement engine.

## Step 1: Core Foundation
> "I need a production-grade Go REST API for a debt-settlement engine. Use Clean Architecture. CRITICAL: Do not use float64 for money; use shopspring/decimal. I want separate layers for Handlers, Services, Algorithms, and Repositories. Use pgxpool for Postgres."

## Step 2: Implementing the Matching Algorithm
> "Write a Greedy matching algorithm in the algorithms package. It should take a map of balances, separate debtors and creditors, sort them by size, and match them up to minimize transfers. Make it deterministic and add unit tests."

## Step 3: Financial Safety & Splitting
> "Add logic to handle equal splits correctly. If a division isn't even (like 10/3), make sure the last person takes the remainder so the sum is always exact. Add validation to reject duplicate people in a split or negative amounts."

## Step 4: Strategy Comparison
> "I want to see how much we are saving with the optimization. Add a second 'Naive' settlement strategy and build a comparison endpoint that tells me the transaction count differences and the percentage gain."

## Step 5: Final Production Audit
> "Do a full audit. Remove any experimental analytics, clean up the routes, and add global panic recovery. Make sure the health check pings the database. Finally, write the documentation in a clear, human-readable way for developers."

## Step 6: Documentation Overhaul
> "Redo the documentation to remove robotic symbols and AI-like formatting. Make it sound like a developer wrote it. Add scenario walkthroughs and a design overview that explains the 'Penny Problem' and float issues."

## Step 7: Localization & Refinement
> "Update the documentation to use the Indian Rupee symbol (₹) instead of the dollar sign ($) across all examples and scenarios. Ensure that the tone remains humanized and all placeholders are professional."

## Step 8: Cross-Platform Diagram Fix
> "The Mermaid diagrams aren't rendering on all live links. Replace the raw Mermaid code and ASCII sketches with a high-compatibility live-rendered image link. Ensure the architecture diagram is visible as a real image in ANY Markdown viewer."
