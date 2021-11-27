<template>
  <div class='column q-gutter-md' style='width: 100%'>

    <q-table
      title='Todos'
      width='100%'
      :columns='columns'
      :rows='todos'
      grid
      :loading='requesting'
      v-model:pagination='pagination'
      :filter='filter'
      @request='onRequest'
    >
      <template v-slot:top-left>
        <q-btn label='タスクを追加' color='secondary' @click='addDialog'></q-btn>
      </template>

      <template v-slot:top-right>
        <q-input borderless dense debounce='300' v-model='filter' placeholder='タスクを検索...'>
          <template v-slot:append>
            <q-icon name='search' />
          </template>
        </q-input>
      </template>

      <template v-slot:item='props'>
        <div class='q-pa-md'>
          <q-card bordered>
            <q-list padding>
              <q-item>
                <q-item-section>
                  {{ props.row.title }}
                  <q-popup-edit :model-value='props.row.title' buttons auto-save
                                @save='(val) => updateTitle(props.row.id, val)' title='タイトルを更新' v-slot='scope'>
                    <q-input v-model='scope.value' dense autofocus counter @keyup.enter='scope.set' />
                  </q-popup-edit>
                </q-item-section>
                <q-item-section side>
                  <q-btn flat round icon='close' color='negative' @click='deleteConfirm(props.row.id)'></q-btn>
                </q-item-section>
              </q-item>

              <q-item>
                <q-item-section>
                  <q-checkbox :model-value='props.row.done'
                              @update:model-value='(val) => updateDone(props.row.id, val)'></q-checkbox>
                </q-item-section>
                <q-item-section>
                  {{ props.row.done ? '完了' : '未完' }}
                </q-item-section>
              </q-item>

              <q-item>
                <q-item-section>
                  <q-item-label>作成日時</q-item-label>
                </q-item-section>
                <q-item-section>
                  {{ new Date(props.row.createdAt).toLocaleString('ja-JP') }}
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>
      </template>
    </q-table>

    <p style='color: lightslategray'>タスクのタイトルをクリックして編集</p>
  </div>
</template>

<script lang='ts'>
import { defineComponent, ref, Ref, onMounted } from 'vue';
import { Todo } from './models';
import { api } from 'boot/axios';
import { QVueGlobals, useQuasar } from 'quasar';

const columns = [
  { name: 'title', label: 'タイトル', field: 'title', required: true, align: 'left' },
  { name: 'done', label: 'ステータス', field: 'done', required: true, align: 'center' },
  { name: 'createdAt', label: '作成日時', field: 'createdAt', required: true, align: 'center' }
];

function useTodoUpdate($q: QVueGlobals, requesting: Ref<boolean>, todos: Ref<Todo[]>) {
  const addTask = (title: string) => {
    requesting.value = true;
    api.post('/api/tasks', { title })
      .then((res) => {
        todos.value.push(res.data);

        $q.notify({
          color: 'positive',
          position: 'bottom',
          message: 'タスクを追加しました。',
          icon: 'done'
        });
      })
      .catch(() => {
        $q.notify({
          color: 'negative',
          position: 'bottom',
          message: 'タスクの追加に失敗しました。',
          icon: 'report_problem'
        });
      })
      .finally(() => {
        requesting.value = false;
      });
  };
  const addDialog = () => {
    $q.dialog({
      title: 'タスクを追加',
      message: 'タスクのタイトルを入力してください。',
      prompt: {
        model: '',
        type: 'text'
      },
      cancel: true
    }).onOk((title: string) => {
      addTask(title);
    });
  };

  const updateTitle = (id: number, title: string) => {
    requesting.value = true;
    api.patch(`/api/tasks/${id}`, { title })
      .then(() => {
        const task = todos.value.find((todo) => todo.id === id);
        if (task) task.title = title;

        $q.notify({
          color: 'positive',
          position: 'bottom',
          message: 'タスクのタイトルを変更しました。',
          icon: 'done'
        });
      })
      .catch(() => {
        $q.notify({
          color: 'negative',
          position: 'bottom',
          message: 'タスクのタイトル変更に失敗しました。',
          icon: 'report_problem'
        });
      })
      .finally(() => {
        requesting.value = false;
      });
  };
  const updateDone = (id: number, done: boolean) => {
    requesting.value = true;
    api.patch(`/api/tasks/${id}`, { done })
      .then(() => {
        const task = todos.value.find((todo) => todo.id === id);
        if (task) task.done = done;

        $q.notify({
          color: 'positive',
          position: 'bottom',
          message: 'タスクの完了状態を変更しました。',
          icon: 'done'
        });
      })
      .catch(() => {
        $q.notify({
          color: 'negative',
          position: 'bottom',
          message: 'タスクの完了状態変更に失敗しました。',
          icon: 'report_problem'
        });
      })
      .finally(() => {
        requesting.value = false;
      });
  };
  const deleteTask = (id: number) => {
    requesting.value = true;
    api.delete(`/api/tasks/${id}`)
      .then(() => {
        todos.value = todos.value.filter((todo) => todo.id !== id);

        $q.notify({
          color: 'positive',
          position: 'bottom',
          message: 'タスクを削除しました。',
          icon: 'done'
        });
      })
      .catch(() => {
        $q.notify({
          color: 'negative',
          position: 'bottom',
          message: 'タスクの削除に失敗しました。',
          icon: 'report_problem'
        });
      })
      .finally(() => {
        requesting.value = false;
      });
  };
  const deleteConfirm = (id: number) => {
    $q.dialog({
      title: 'タスクを削除',
      message: '本当にこのタスクを削除しますか？',
      cancel: true
    }).onOk(() => deleteTask(id));
  };
  return { addDialog, updateTitle, updateDone, deleteConfirm };
}

export default defineComponent({
  name: 'TodoList',

  props: {},

  setup() {
    const $q = useQuasar();

    const todos = ref([] as Todo[]);
    const filter = ref('');
    const requesting = ref(false);
    const pagination = ref({
      descending: false,
      page: 1,
      rowsPerPage: 5,
      rowsNumber: 10
    });

    function onRequest(props: { pagination: { page: number, rowsPerPage: number }; filter: string }) {
      const { page, rowsPerPage } = props.pagination
      const filter = props.filter

      requesting.value = true;

      const params: Record<string, string | number> = {}
      if (filter) {
        params.title = filter
      }
      // TODO: done task only
      if (rowsPerPage > 0) {
        params.limit = rowsPerPage
        params.offset = rowsPerPage * (page - 1)
      }
      api.get('/api/tasks', { params })
      .then((res) => {
        const data = res.data as { count: number; tasks: Todo[] }

        todos.value = data.tasks
        pagination.value.rowsNumber = data.count
        pagination.value.page = page
        pagination.value.rowsPerPage = rowsPerPage
      })
      .catch(() => {
        $q.notify({
          color: 'negative',
          position: 'bottom',
          message: 'タスクリストの取得に失敗しました。',
          icon: 'report_problem'
        })
      })
      .finally(() => {
        requesting.value = false;
      })
    }

    onMounted(() => {
      onRequest({
        pagination: pagination.value,
        filter: ''
      });
    });

    return { todos, columns, requesting, filter, pagination, onRequest, ...useTodoUpdate($q, requesting, todos) };
  }
});
</script>
