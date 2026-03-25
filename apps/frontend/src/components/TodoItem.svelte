<script>
  import { lists } from '../stores/lists.js';

  export let todo;
  export let listId;

  let editing = false;
  let editText = todo.text;

  function saveEdit() {
    if (editText.trim()) lists.updateTodo(listId, todo.id, editText);
    editing = false;
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') saveEdit();
    if (e.key === 'Escape') { editing = false; editText = todo.text; }
  }
</script>

<div class="todo" class:done={todo.done}>
  <button class="check" on:click={() => lists.toggleTodo(listId, todo.id)}>
    {#if todo.done}
      <span class="chk done">[X]</span>
    {:else}
      <span class="chk">[ ]</span>
    {/if}
  </button>

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

  <button class="del" on:click={() => lists.deleteTodo(listId, todo.id)}>[DEL]</button>
</div>

<style>
  .todo {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.35rem 0;
    border-bottom: 1px solid var(--border-faint);
  }
  .todo:last-child { border-bottom: none; }
  .todo.done { opacity: 0.4; }

  .check {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    flex-shrink: 0;
  }
  .chk { font-family: var(--font-pixel); font-size: 0.45rem; color: var(--muted); }
  .chk.done { color: var(--accent); text-shadow: 0 0 4px var(--accent-glow); }

  .todo-text {
    flex: 1;
    font-size: 0.45rem;
    color: var(--text-secondary);
    cursor: pointer;
    line-height: 1.6;
    word-break: break-word;
  }
  .todo.done .todo-text {
    text-decoration: line-through;
    color: var(--muted);
  }

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

  .del {
    background: none;
    border: none;
    cursor: pointer;
    color: transparent;
    font-family: var(--font-pixel);
    font-size: 0.4rem;
    padding: 0.1rem 0.2rem;
    transition: color 0.1s;
    flex-shrink: 0;
  }
  .todo:hover .del { color: var(--muted); }
  .del:hover { color: var(--danger) !important; }
</style>