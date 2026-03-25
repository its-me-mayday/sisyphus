<script>
  import { lists } from '../stores/lists.js';

  export let todo;
  export let listId;

  let editing = false;
  let editText = todo.text;

  const PRIORITY = {
    low:    { label: 'LOW',  color: 'var(--muted)' },
    medium: { label: 'MED',  color: 'var(--amber)' },
    high:   { label: 'HIGH', color: 'var(--danger)' },
  };

  const PRIORITY_CYCLE = { low: 'medium', medium: 'high', high: 'low' };

  function saveEdit() {
    if (editText.trim()) lists.updateTodo(listId, todo.id, editText);
    editing = false;
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') saveEdit();
    if (e.key === 'Escape') { editing = false; editText = todo.text; }
  }

  function cyclePriority() {
    lists.setTodoPriority(listId, todo.id, PRIORITY_CYCLE[todo.priority ?? 'medium']);
  }

  function handleDueDate(e) {
    lists.setTodoDueDate(listId, todo.id, e.target.value || null);
  }

  $: p = PRIORITY[todo.priority ?? 'medium'];

  $: isOverdue = todo.dueDate && !todo.done && new Date(todo.dueDate) < new Date();
</script>

<div class="todo" class:done={todo.done} class:overdue={isOverdue}>
  <!-- Toggle -->
  <button class="check" on:click={() => lists.toggleTodo(listId, todo.id)}>
    {#if todo.done}
      <span class="chk done">[X]</span>
    {:else}
      <span class="chk">[ ]</span>
    {/if}
  </button>

  <!-- Priority badge -->
  <button class="priority" style="color: {p.color};" on:click={cyclePriority} title="cambia priorità">
    {p.label}
  </button>

  <!-- Text -->
  {#if editing}
    <input
      class="edit-input"
      bind:value={editText}
      on:blur={saveEdit}
      on:keydown={handleKeydown}
    />
  {:else}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <span class="todo-text" on:dblclick={() => { editing = true; editText = todo.text; }}>
      {todo.text}
    </span>
  {/if}

  <!-- Due date -->
  <input
    class="due-date"
    type="date"
    value={todo.dueDate ?? ''}
    on:change={handleDueDate}
    title="scadenza"
  />

  <!-- Delete -->
  <button class="del" on:click={() => lists.deleteTodo(listId, todo.id)}>[DEL]</button>
</div>

<style>
  .todo {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.35rem 0;
    border-bottom: 1px solid var(--border-faint);
  }
  .todo:last-child { border-bottom: none; }
  .todo.done    { opacity: 0.4; }
  .todo.overdue { border-left: 2px solid var(--danger); padding-left: 0.4rem; }

  .check {
    background: none; border: none; cursor: pointer;
    padding: 0; flex-shrink: 0;
  }
  .chk { font-family: var(--font-pixel); font-size: 0.45rem; color: var(--muted); }
  .chk.done { color: var(--accent); text-shadow: 0 0 4px var(--accent-glow); }

  .priority {
    background: none;
    border: 1px solid currentColor;
    font-family: var(--font-pixel);
    font-size: 0.38rem;
    cursor: pointer;
    padding: 0.1rem 0.3rem;
    flex-shrink: 0;
    opacity: 0.8;
    transition: opacity 0.1s;
  }
  .priority:hover { opacity: 1; }
  .todo.done .priority { opacity: 0.3; }

  .todo-text {
    flex: 1;
    font-size: 0.45rem;
    color: var(--text-secondary);
    cursor: pointer;
    line-height: 1.6;
    word-break: break-word;
  }
  .todo.done .todo-text { text-decoration: line-through; color: var(--muted); }

  .edit-input {
    flex: 1;
    background: var(--accent-subtle);
    border: 1px solid var(--accent);
    color: var(--text-primary);
    font-family: var(--font-pixel);
    font-size: 0.45rem;
    padding: 0.15rem 0.3rem;
    outline: none;
  }

  .due-date {
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--border);
    color: var(--muted);
    font-family: var(--font-pixel);
    font-size: 0.35rem;
    padding: 0.1rem;
    cursor: pointer;
    flex-shrink: 0;
    width: 90px;
    outline: none;
    color-scheme: dark;
  }
  .due-date:hover { border-color: var(--border-bright); color: var(--text-secondary); }
  .todo.overdue .due-date { color: var(--danger); border-color: var(--danger); }

  .del {
    background: none; border: none; cursor: pointer;
    color: transparent;
    font-family: var(--font-pixel); font-size: 0.4rem;
    padding: 0.1rem 0.2rem;
    transition: color 0.1s; flex-shrink: 0;
  }
  .todo:hover .del { color: var(--muted); }
  .del:hover { color: var(--danger) !important; }
</style>