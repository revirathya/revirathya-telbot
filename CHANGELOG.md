# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Pre-release - Command, Conversation, Context and gRPC contracts may change without notice.

## [Unreleased]

## [0.0.1] - 2026-09-09

Initial pre-release.

### Command: `/expense`

Expenses related conversation, gateway for expenses conversation: [Record Expense](#conversation-record-expense), [List Expenses](#conversation-list-expenses), [Summary Expenses](#conversation-summary-expenses)

**Flow:**

### Command: `/cancel`

Force to exit conversation

### Conversation: Record Expense

Recording expense

**Flow:**

### Conversation: List Expenses

List recorded expenses

**Flow:**

### Conversation: Summary Expenses

Generate summary of pre-defined reports regarding expenses

**Flow:**

### Subflow
- `Numpad`: callback-based number inputs
- `Date Picker`: callback-based calendar navigator and date-picker