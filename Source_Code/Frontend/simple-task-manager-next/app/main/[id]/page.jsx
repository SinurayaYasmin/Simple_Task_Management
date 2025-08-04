'use client';

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import axios from "axios";

const DashboardPage = () => {
  const { id } = useParams();
  const [user, setUser] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [newTask, setNewTask] = useState({
    title: "",
    description: "",
    status: "Todo",
    deadline: "",
  });

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await axios.get(`http://localhost:8080/users/${id}`);
        setUser(res.data);
      } catch (err) {
        console.error("Failed to fetch user:", err);
      }
    };

    const fetchTasks = async () => {
      try {
        const res = await axios.get(`http://localhost:8080/tasks/getAllTask`);
        setTasks(res.data);
      } catch (err) {
        console.error("Failed to fetch tasks:", err);
      }
    };

    if (id) {
      fetchUser();
      fetchTasks();
    }
  }, [id]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewTask((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleAddTask = async () => {
    try {
      const res = await axios.post(`http://localhost:8080/tasks/${id}/createTask`, {
        title: newTask.title,
        description: newTask.description,
        status: newTask.status,
        deadline: new Date(newTask.deadline).toISOString(),
      });
      setTasks([...tasks, res.data]);
      setNewTask({
        title: "",
        description: "",
        status: "Todo",
        deadline: "",
      });
    } catch (err) {
      console.error("Failed to add task:", err);
    }
  };

  const handleDeleteTask = async (taskId) => {
    try {
      await axios.delete(`http://localhost:8080/tasks/deleteTask/${taskId}`);
      setTasks(tasks.filter((task) => task.id !== taskId));
    } catch (err) {
      console.error("Failed to delete task:", err);
    }
  };

  if (!user) return <div>Loading user data...</div>;

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Welcome, {user.username}!</h1>

      <div className="mb-6 space-y-3 bg-gray-50 p-4 rounded shadow">
        <h2 className="text-xl font-semibold">Create New Task</h2>
        <input
          name="title"
          value={newTask.title}
          onChange={handleInputChange}
          placeholder="Title"
          className="w-full border p-2 rounded"
        />
        <textarea
          name="description"
          value={newTask.description}
          onChange={handleInputChange}
          placeholder="Description"
          className="w-full border p-2 rounded"
        />
        <select
          name="status"
          value={newTask.status}
          onChange={handleInputChange}
          className="w-full border p-2 rounded"
        >
          <option value="Todo">Todo</option>
          <option value="In Progress">In Progress</option>
          <option value="Done">Done</option>
        </select>
        <input
          name="deadline"
          type="datetime-local"
          value={newTask.deadline}
          onChange={handleInputChange}
          className="w-full border p-2 rounded"
        />
        <button
          onClick={handleAddTask}
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Add Task
        </button>
      </div>

      <div>
        <h2 className="text-xl font-semibold mb-2">Task List</h2>
        <ul className="space-y-3">
          {tasks.map((task) => (
            <li
              key={task.id}
              className="flex justify-between items-center bg-white p-4 rounded shadow"
            >
              <div>
                <p className="font-bold">{task.title}</p>
                <p className="text-sm text-gray-600">{task.description}</p>
                <p className="text-xs text-gray-400">
                  Status: {task.status} | Deadline: {new Date(task.deadline).toLocaleString()}
                </p>
              </div>
              <button
                onClick={() => handleDeleteTask(task.id)}
                className="bg-red-500 text-white px-3 py-1 rounded"
              >
                Delete
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default DashboardPage;
