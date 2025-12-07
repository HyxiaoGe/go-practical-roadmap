// 任务系统前端JavaScript代码

// 全局变量
let ws = null;
let tasks = {};
let isConnected = false;

// WebSocket连接管理
class WebSocketManager {
    constructor(url) {
        this.url = url;
        this.reconnectInterval = 5000;
        this.maxReconnectAttempts = 10;
        this.reconnectAttempts = 0;
    }

    connect() {
        try {
            this.ws = new WebSocket(this.url);
            this.setupEventHandlers();
        } catch (error) {
            console.error('Failed to create WebSocket:', error);
            this.handleReconnect();
        }
    }

    setupEventHandlers() {
        this.ws.onopen = (event) => {
            console.log('WebSocket connected');
            isConnected = true;
            this.reconnectAttempts = 0;
            this.updateConnectionStatus(true);
            this.addToLog('WebSocket连接已建立');
        };

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleMessage(data);
            } catch (error) {
                console.error('Failed to parse message:', error);
            }
        };

        this.ws.onclose = (event) => {
            console.log('WebSocket disconnected');
            isConnected = false;
            this.updateConnectionStatus(false);
            this.addToLog('WebSocket连接已断开');
            this.handleReconnect();
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.addToLog('WebSocket连接错误: ' + error);
        };
    }

    handleMessage(data) {
        if (data.type === 'connection') {
            this.addToLog('连接成功: ' + data.message);
        } else {
            // 更新任务状态
            updateTask(data);
            this.addToLog(`任务 ${data.task_id} 状态更新为 ${data.status}`);
        }
    }

    updateConnectionStatus(connected) {
        const statusElement = document.getElementById('ws-status');
        if (connected) {
            statusElement.textContent = '已连接';
            statusElement.className = 'websocket-status ws-connected';
        } else {
            statusElement.textContent = '已断开';
            statusElement.className = 'websocket-status ws-disconnected';
        }
    }

    handleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            this.addToLog(`尝试重新连接 (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

            setTimeout(() => {
                this.connect();
            }, this.reconnectInterval);
        } else {
            this.addToLog('达到最大重连次数，停止重连');
        }
    }

    send(data) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        }
    }

    close() {
        if (this.ws) {
            this.ws.close();
        }
    }

    addToLog(message) {
        const logContent = document.getElementById('log-content');
        const timestamp = new Date().toLocaleTimeString();
        logContent.textContent += `[${timestamp}] ${message}\n`;
        logContent.scrollTop = logContent.scrollHeight;
    }
}

// 任务管理功能
class TaskManager {
    static submitTask() {
        const taskName = document.getElementById('task-name').value;
        const taskPayload = document.getElementById('task-payload').value;

        let payload;
        try {
            payload = JSON.parse(taskPayload);
        } catch (e) {
            alert('无效的JSON格式');
            return;
        }

        const requestData = {
            name: taskName,
            payload: payload
        };

        fetch('/api/v1/tasks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.json())
        .then(data => {
            wsManager.addToLog(`任务提交成功: ${data.id}`);
            // 添加到本地任务列表
            tasks[data.id] = {
                task_id: data.id,
                status: 'pending',
                progress: 0,
                timestamp: new Date().toISOString()
            };
            updateTaskList();
            updateStats();
        })
        .catch(error => {
            wsManager.addToLog(`任务提交失败: ${error}`);
            console.error('Error:', error);
        });
    }

    static cancelTask(taskId) {
        fetch(`/api/v1/tasks/${taskId}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                wsManager.addToLog(`任务 ${taskId} 已取消`);
            } else {
                wsManager.addToLog(`取消任务 ${taskId} 失败`);
            }
        })
        .catch(error => {
            wsManager.addToLog(`取消任务 ${taskId} 错误: ${error}`);
            console.error('Error:', error);
        });
    }

    static viewTaskDetails(taskId) {
        const task = tasks[taskId];
        if (task) {
            alert(`任务详情:\nID: ${task.task_id}\n状态: ${task.status}\n进度: ${task.progress || 0}%\n时间戳: ${task.timestamp}`);
        }
    }
}

// UI更新功能
function updateTask(taskData) {
    // 更新任务数据
    tasks[taskData.task_id] = taskData;

    // 更新UI
    updateTaskList();
    updateStats();
    renderTaskItem(taskData);
}

function updateTaskList() {
    const taskListElement = document.getElementById('task-list');

    // 清空现有内容
    taskListElement.innerHTML = '';

    // 按时间倒序排列任务
    const taskArray = Object.values(tasks).sort((a, b) =>
        new Date(b.timestamp) - new Date(a.timestamp)
    );

    // 显示最新的20个任务
    taskArray.slice(0, 20).forEach(task => {
        renderTaskItem(task);
    });
}

function renderTaskItem(task) {
    const taskListElement = document.getElementById('task-list');
    const existingElement = document.getElementById(`task-${task.task_id}`);

    const taskHtml = `
        <div class="task-item" id="task-${task.task_id}">
            <div class="task-header">
                <span class="task-id">${task.task_id}</span>
                <span class="task-status status-${task.status}">${task.status}</span>
            </div>
            <div>进度: ${task.progress || 0}%</div>
            <div class="task-progress">
                <div class="progress-bar" style="width: ${task.progress || 0}%"></div>
            </div>
            ${task.error ? `<div style="color: red;">错误: ${task.error}</div>` : ''}
            <div class="task-actions">
                <button onclick="TaskManager.cancelTask('${task.task_id}')">取消任务</button>
                <button onclick="TaskManager.viewTaskDetails('${task.task_id}')">查看详情</button>
            </div>
        </div>
    `;

    if (existingElement) {
        existingElement.outerHTML = taskHtml;
    } else {
        taskListElement.insertAdjacentHTML('afterbegin', taskHtml);
    }
}

function updateStats() {
    const stats = {
        total: 0,
        pending: 0,
        running: 0,
        completed: 0,
        failed: 0,
        cancelled: 0
    };

    Object.values(tasks).forEach(task => {
        stats.total++;
        stats[task.status]++;
    });

    document.getElementById('total-tasks').textContent = stats.total;
    document.getElementById('running-tasks').textContent = stats.running;
    document.getElementById('completed-tasks').textContent = stats.completed;
    document.getElementById('failed-tasks').textContent = stats.failed;
}

// 页面初始化
document.addEventListener('DOMContentLoaded', function() {
    // 初始化WebSocket管理器
    const wsUrl = `ws://${window.location.host}/api/v1/ws/tasks`;
    wsManager = new WebSocketManager(wsUrl);
    wsManager.connect();

    // 定期刷新任务列表
    setInterval(function() {
        fetch('/api/v1/tasks')
        .then(response => response.json())
        .then(data => {
            // 更新任务数据
            data.tasks.forEach(task => {
                tasks[task.id] = {
                    task_id: task.id,
                    status: task.status,
                    progress: task.progress,
                    timestamp: task.created_at
                };
            });
            updateTaskList();
            updateStats();
        })
        .catch(error => {
            console.error('Error fetching tasks:', error);
        });
    }, 5000);

    // 绑定提交按钮事件
    document.getElementById('submit-task-btn').addEventListener('click', TaskManager.submitTask);
});