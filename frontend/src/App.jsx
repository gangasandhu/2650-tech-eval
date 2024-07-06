import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [todos, setTodos] = useState([]);
  const [title, setTitle] = useState('');

  useEffect(() => {
    const fetchTodos = async () => {
      try {
        const response = await fetch('http://localhost:3000/todos');
        const data = await response.json();
        setTodos(data);
      } catch (error) {
        console.error('Error fetching todos:', error);
      }
    };

    fetchTodos();
  }, []);

  const addTodo = async () => {
    try {
      const response = await fetch('http://localhost:3000/todo', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title, done: false }),
      });
      const newTodo = await response.json();
      setTodos([...todos, newTodo]);
      setTitle('');
    } catch (error) {
      console.error('Error adding todo:', error);
    }
  };

  const deleteTodo = async (id) => {
    try {
      await fetch(`http://localhost:3000/todo/${id}`, { method: 'DELETE' });
      setTodos(todos.filter(todo => todo.id !== id));
    } catch (error) {
      console.error('Error deleting todo:', error);
    }
  };

  const toggleTodoStatus = async (id, done) => {
    try {
      const response = await fetch(`http://localhost:3000/todo/${id}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ done: !done }),
      });
      const updatedTodo = await response.json();
      setTodos(todos.map(todo => (todo.id === id ? {...todo, done: !done}  : todo)));
    } catch (error) {
      console.error('Error updating todo status:', error);
    }
  };

  const updateTodoTitle = async (id, title) => {
    try {
      const response = await fetch(`http://localhost:3000/todo/${id}/title`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title }),
      });
      const updatedTodo = await response.json();
      setTodos(todos.map(todo => (todo.id === id ? updatedTodo : todo)));
    } catch (error) {
      console.error('Error updating todo title:', error);
    }
  };

  const handleTitleChange = (id, title) => {
    setTodos(todos.map(todo => (todo.id === id ? { ...todo, title } : todo)));
  };


  return (
    <div className="App">
      <h1>Todo List</h1>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Enter todo"
      />
      <button onClick={addTodo}>Add Todo</button>
      <ul>
        {todos.map(todo => (
          <li key={todo.id}>
            <input
              type="checkbox"
              checked={todo.done}
              onChange={() => toggleTodoStatus(todo.id, todo.done)}
            />
           <input
              type="text"
              value={todo.title}
              onChange={(e) => handleTitleChange(todo.id, e.target.value)}
              onBlur={() => updateTodoTitle(todo.id, todo.title)}
              style={{ textDecoration: todo.done ? 'line-through' : 'none' }}
            />
            <button onClick={() => deleteTodo(todo.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
