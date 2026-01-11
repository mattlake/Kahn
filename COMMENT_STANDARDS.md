# Kahn Task Manager - Comment Standards

## 🚫 When NOT to Add Comments

### 1. Redundant Function Descriptions
**AVOID:** Comments that simply repeat the function name
```go
// ❌ BAD: Add adds an item to the repository
func (r *Repo) Add(item Item) error

// ❌ BAD: Delete removes a task from the repository  
func (r *Repo) Delete(id string) error
```

**REASON:** The function name already describes what it does.

### 2. Obvious Code Explanations
**AVOID:** Comments that state what the code obviously does
```go
// ❌ BAD: Use domain validation for data integrity
if err := task.Validate(); err != nil {
    return nil, err
}

// ❌ BAD: Set the task priority
task.Priority = priority

// ❌ BAD: Return the result
return result
```

**REASON:** Good code should be self-documenting.

### 3. Simple Assignment Comments
**AVOID:** Comments that explain basic variable assignments
```go
// ❌ BAD: Initialize the active index
activeIndex := 0

// ❌ BAD: Create a new task
task := NewTask(name, desc, projectID)
```

**REASON:** Variable names and function calls make the purpose clear.

### 4. Test Case Explanations
**AVOID:** Comments that repeat test case names or obvious setup
```go
// ❌ BAD: Test adding task to project
project.AddTask(*task)

// ❌ BAD: Update status to in progress
task.Status = InProgress

// ❌ BAD: Multiple calls to ensure uniqueness
id1 := generateID()
id2 := generateID()
```

**REASON:** Test names should describe the scenario; code should be clear.

## ✅ When to Add Comments

### 1. Complex Business Logic
**USE:** Comments explaining non-obvious business rules
```go
// ✅ GOOD: Different ordering based on status - Not Started: priority DESC, then created_at ASC (oldest highest priority first)
if status == domain.NotStarted {
    // Sort by priority descending, then creation time ascending
} else {
    // Sort by updated_at descending (newest changes first)
}
```

### 2. Security and Safety Considerations
**USE:** Comments explaining security measures
```go
// ✅ GOOD: Clean path to resolve any .. sequences for security
dbPath := filepath.Clean(inputPath)

// ✅ GOOD: Validate database path to prevent directory traversal
if !strings.Contains(dbPath, "..") {
    return fmt.Errorf("invalid path")
}
```

### 3. Performance Optimizations
**USE:** Comments explaining performance decisions
```go
// ✅ GOOD: PERFORMANCE: Use cached style objects to avoid allocations
if isActiveList && isSelected {
    return selectedStyle.Render(title)
}
```

### 4. Complex Algorithm Explanations
**USE:** Comments explaining algorithmic choices
```go
// ✅ GOOD: Use binary search for O(log n) performance instead of O(n) linear scan
func findItem(items []Item, target string) int {
    // Implementation...
}

// ✅ GOOD: Preserve existing selection states during list updates to maintain UI consistency
// This prevents jarring jumps in cursor position when data changes
selections := saveSelectionStates()
```

### 5. Workarounds and Temporary Solutions
**USE:** Comments explaining why suboptimal code exists
```go
// ✅ GOOD: TODO: Replace with proper caching layer when available
// Temporary workaround for database connection pooling issues
func getConnection() *sql.DB {
    // Workaround implementation...
}

// ✅ GOOD: Workaround for lipgloss limitation with nested styling
// Remove when library supports border styling on containers
func customBorder() string {
    // Workaround code...
}
```

## 📏�️ Code Should Be Self-Documenting

### 1. Clear Variable Names
```go
// ✅ GOOD: maxRetries, connectionTimeout, isActive
// ❌ BAD: mr, ct, flag
```

### 2. Clear Function Names
```go
// ✅ GOOD: ValidateUserCredentials, CalculateTotalPrice, GetActiveProject
// ❌ BAD: Validate, Calculate, GetProject
```

### 3. Consistent Structure
```go
// ✅ GOOD: Well-structured code needs fewer comments
func (ts *TaskService) CreateTask(name, description, projectID string) (*Task, error) {
    // Input validation
    if err := validateInputs(name, projectID); err != nil {
        return nil, err
    }
    
    // Business logic
    task := domain.NewTask(name, description, projectID)
    
    // Persistence
    if err := ts.taskRepo.Create(task); err != nil {
        return nil, err
    }
    
    return task, nil
}
```

## 🎯 Principle: "Code as Documentation"

The best comments are no comments at all. Write code that is so clear that comments become redundant.

**Remember:** Every comment you add is technical debt that must be maintained alongside the code. If the comment and code diverge, the comment becomes a source of confusion.

## 📝 Review Process

1. **Code Review:** Check for needless comments during code review
2. **Refactor:** Improve code clarity instead of adding comments  
3. **Question:** Ask "Would a new developer understand this without the comment?"
4. **Remove:** Delete comments that become obsolete after refactoring
5. **Audit:** Regularly review existing comments for ongoing relevance

---

**Last Updated:** January 11, 2026  
**Enforced:** All new code should follow these standards