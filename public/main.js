(function (w) {
    var PAGE_SIZE = 3;
    w.fn = {};
    w.instance = {};

    function Modal() {
        if (!(this instanceof Modal)) {
            return new Modal();
        }
        this.modalContainer = $('.modal-container');
    };

    Modal.prototype.init = function () {
        var self = this;
        $('.modal__ok-button').on('click', function (_) {
            self.submitNewTodo();
        });

        $('.modal__cancel-button').on('click', function (_) {
            $('.modal__input').val('');
            self.hide();
        });

        $('.modal__input').on('keyup', function (e) {
            e.key.toLowerCase() === 'enter' && self.submitNewTodo();
        });
    };

    Modal.prototype.getContent = function () {
        var content  = $('.modal__input').val();
        return content;
    };

    Modal.prototype.submitNewTodo = function () {
        var content = this.getContent();
        var data = {
            content: content,
            rootId: -1,
            parentId: []
        };
        Service.createTodo(data).then(function (ok) {
            if (ok) {
                document.location.reload();
            } else {
                alert('잠시 후 다시 시도해 주세요.');
            }
        });
    };

    Modal.prototype.show = function () {
        this.modalContainer.addClass('show');
    };

    Modal.prototype.hide = function () {
        this.modalContainer.removeClass('show');
    };
    w.instance.modal = new Modal();

    function AddTodoButton(modal) {
        if (!(this instanceof AddTodoButton)) {
            return new AddTodoButton();
        }
        this.modal = modal;
    };

    AddTodoButton.prototype.init = function () {
        var self = this;
        $('.todolist__add-button').on('click', function (_) {
            self.modal.show();
        });
    };
    w.instance.addTodoButton = new AddTodoButton(w.instance.modal);

    function TodoList() {
        if (!(this instanceof TodoList)) {
            return new TodoList();
        }
        this.dom = $('.todolist__container');
    }

    TodoList.prototype.render = function(todos) {
        var tbody = this.dom.find('tbody');
        tbody.empty();
        todos.map(todo => Todo.toComponent(todo, tbody)).forEach(Todo.render);
    };
    w.instance.todoList = new TodoList();

    function Todo(data, container) {
        if (!(this instanceof Todo)) {
            return new Todo(arguments);
        }
        this.container = container;
        this.id = data.id;
        this.content = data.content;
        this.createdAt = data.createdAt;
        this.modifiedAt = data.modifiedAt;
        this.isComplete = data.isComplete;
        this.rootId = data.rootId;
        this.parentIds = data.parentIds;

        this.dom = this.createDOM();
    }

    Todo.prototype.createDOM = function () {
        var self = this;
        var addChildButton = $('<span>')
            .addClass('item__add-child-button')
            .text('+')
            .on('click', function (_) {
                var emptyData = {
                    id: 0,
                    content: '',
                    createdAt: '',
                    modifiedAt: '',
                    isComplete: false,
                    rootId: self.rootId < 0 ? self.id : self.rootId,
                    parentIds: [ self.id ]
                };
                var newTodo = new Todo(emptyData, self.container);
                newTodo.dom.insertAfter(self.dom);
            });
        var addChild = $('<td>').addClass('item__add-child').append(addChildButton);
        var id = $('<td>').addClass('item__id').text(this.id);
        var contentInput = $('<input>')
            .attr('type', 'text')
            .val(this.content)
            .on('keyup', function (e) {
                if (e.key.toLowerCase() !== 'enter') return true;
                var data = self.getData();
                var updatedData = Object.assign(data, { content: this.value });
                if (self.id === 0) {
                    self.submitNewChild(updatedData).then(function (ok) {
                        if (ok) {
                            document.location.reload();
                        } else {
                            alert('잠시 후 다시 시도해 주세요.');
                        }
                    });
                } else {
                    self.update(updatedData).then(function (ok) {
                        if (!ok) {
                            $(this).prop('checked', !isComplete);
                            alert('잠시 후 다시 시도해 주세요.');
                        }
                    });
                }
            });
        var content = $('<td>')
            .addClass('item__content')
            .append(contentInput);
        var parents = $('<td>').addClass('item__references');
        if (this.parentIds && this.parentIds.length) {
            parents.text(this.parentIds.map(function(pid) {
                return '@' + pid;
            }).join(' '));
        }
        var createdAt = $('<td>').addClass('item__created').text(this.createdAt);
        var modifiedAt = $('<td>').addClass('item__modified').text(this.modifiedAt);
        var checkbox = $('<input>')
            .attr('type', 'checkbox')
            .prop('checked', this.isComplete)
            .on('change', function (e) {
                var isComplete = e.target.checked;
                var data = self.getData();
                var updatedData = Object.assign(data, { isComplete: isComplete });
                self.update(updatedData).then(function (ok) {
                    if (!ok) {
                        alert('업데이트 할 수 없습니다.');
                        document.location.reload();
                    }
                });
            });
        var isComplete = $('<td>').addClass('item__checkbox').append(checkbox);
        var root = $('<tr>')
            .addClass('todolist__item')
            .data('id', this.id)
            .append(addChild, id, content, parents, createdAt, modifiedAt, isComplete);
        return root;
    };

    Todo.prototype.appendToContainer = function () {
        this.container.append(this.dom);
    };

    Todo.prototype.getData = function () {
        var children = this.dom.children();
        var id = parseInt(children[1].innerText);
        var content = $(children[2]).children('input').val();
        var parentIdsString = children[3].innerText.trim();
        var parentIds = parentIdsString.length ? parentIdsString.split(' ').map(function (text) {
            return parseInt(text.replace('@', ''));
        }) : [];
        var createdAt = children[4].innerText;
        var modifiedAt = children[5].innerText;
        var isComplete = $(children[6]).children('input').prop('checked');
        return {
            id: id,
            content: content,
            parentIds: parentIds,
            createdAt: createdAt,
            modifiedAt: modifiedAt,
            isComplete: isComplete
        };
    };
    
    Todo.prototype.update = function (data) {
        return Service.updateTodo(data);
    };

    Todo.prototype.submitNewChild = function (data) {
        return Service.createTodo(data);
    }

    Todo.toComponent = function (data, container) {
        return new Todo(data, container);
    };

    Todo.render = function (component) {
        component.appendToContainer();
    } ;
    w.Todo = Todo;

    function Pagination() {
        this.container = $('.todolist__pagination-container');
    };

    Pagination.prototype.createDOM = function (currentPageNo, totalCount) {
        var startPageNo = this.getStartPageNo(currentPageNo);
        var endPageNo = this.getEndPageNo(startPageNo, totalCount);
        var lis = [];
        for(var idx = startPageNo; idx <= endPageNo; idx++) {
            (function (memory) {
                var anchor = $('<a>')
                .attr('href', '#')
                .text(memory);
                var li = $('<li>')
                    .addClass(currentPageNo === memory ? 'todolist__pagination-button--selected' : 'todolist__pagination-button')
                    .on('click', function (_) {
                        Service.getTodos(memory).then(w.fn.rerender);
                    })
                    .append(anchor);
                lis.push(li);
            })(idx);
        }
        return lis;
    };

    Pagination.prototype.getStartPageNo = function (currentPageNo) {
        var startPageNo = Math.floor(currentPageNo / PAGE_SIZE) + (currentPageNo % PAGE_SIZE === 0 ? 0 : 1);
        return startPageNo;
    };

    Pagination.prototype.getEndPageNo = function (startPageNo, totalCount) {
        var endPageNo = startPageNo + (PAGE_SIZE - 1);
        var veryLastPageNo = Math.floor(totalCount / PAGE_SIZE) + (totalCount % PAGE_SIZE === 0 ? 0 : 1);
        if (endPageNo > veryLastPageNo) {
            endPageNo = veryLastPageNo;
        }
        return endPageNo;
    };

    Pagination.prototype.appendToContainer = function () {
        this.container.empty();
        this.container.append(this.dom);
    };

    Pagination.prototype.render = function (currentPageNo, totalCount) {
        this.dom = this.createDOM(currentPageNo, totalCount);
        this.appendToContainer();
    };
    w.instance.pagination = new Pagination();

    w.fn.rerender = function (data) {
        var todos = data.todos;
        var totalCount = data.totalCount;
        var pageNo = data.pageNo;
        w.instance.todoList.render(todos);
        w.instance.pagination.render(pageNo, totalCount);
    };
})(window, 'Component');

(function (w) {
    var Service = {};
    Service.createTodo = function (todo) {
        return fetch('api/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(todo)
        }).then(function (response) {
            return response.ok;
        }).catch(function (err) {
            console.dir(err);
            return false;
        });
    };

    Service.getTodos = function (pageNo) {
        pageNo = pageNo || 1;
        return fetch('api/todos?pageNo=' + pageNo, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
            }
        }).then(function (response) {
            return response.json();
        }).then(function (json) {
            return Object.assign(json, { pageNo: pageNo });
        })
    };

    Service.updateTodo = function (todo) {
        return fetch('api/todos', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(todo)
        }).then(function (response) {
            return response.ok;
        }).catch(function (err) {
            console.dir(err);
            return false;
        });
    };

    w.Service = Service;
})(window, 'Service');

(function () {
    instance.modal.init();
    instance.addTodoButton.init();
    Service.getTodos().then(fn.rerender);
})();