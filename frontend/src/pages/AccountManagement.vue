<template>
  <q-page padding class='column'>

    <h6>アカウント管理</h6>

    <div v-if='!loading'>
      <p>名前: {{ data.name }}</p>
      <p>作成日時: {{ new Date(data.createdAt).toLocaleString('ja-JP') }}</p>
      <p>情報最終更新日時: {{ new Date(data.updatedAt).toLocaleString('ja-JP') }}</p>
    </div>
    <div v-else>
      <p>Loading...</p>
    </div>

    <div class='q-gutter-md row'>
      <q-btn label='名前を変更する' color='secondary' @click='changeNamePrompt'></q-btn>
      <q-btn label='アカウントを削除する' color='negative' @click='deleteConfirm($router)'></q-btn>
    </div>

  </q-page>
</template>

<script lang='ts'>
import { defineComponent, Ref, ref } from 'vue';
import { api } from 'boot/axios';
import { useQuasar } from 'quasar';
import { AxiosError } from 'axios';
import { Router } from 'vue-router';

interface AccountData {
  name: string;
  createdAt: string;
  updatedAt: string;
}

export default defineComponent({
  name: 'AccountManagement',

  components: {},

  setup() {
    const $q = useQuasar();

    const data: Ref<AccountData> = ref({ name: '', createdAt: '', updatedAt: '' });
    const loading = ref(false);

    const reload = () => {
      loading.value = true;
      api.get('/api/users/me')
        .then((res) => {
          data.value = res.data as AccountData;
        })
        .catch(() => {
          $q.notify({
            color: 'negative',
            position: 'bottom',
            message: 'アカウント情報の取得に失敗しました。',
            icon: 'report_problem'
          });
        })
        .finally(() => {
          loading.value = false;
        });
    };

    const changeName = (name: string) => {
      api.patch('/api/users/me', { name })
        .then(() => {
          $q.notify({
            color: 'positive',
            position: 'bottom',
            message: '名前を変更しました。',
            icon: 'done'
          });
        })
        .catch((res) => {
          if ((res as AxiosError).response?.status === 409) {
            $q.notify({
              color: 'negative',
              position: 'bottom',
              message: 'すでに使われているユーザー名です。',
              icon: 'report_problem'
            });
          } else {
            $q.notify({
              color: 'negative',
              position: 'bottom',
              message: '名前の変更に失敗しました。',
              icon: 'report_problem'
            });
          }
        })
        .finally(() => {
          reload();
        });
    };
    const changeNamePrompt = () => {
      $q.dialog({
        title: '名前を変更する',
        message: '新しいユーザー名を入力してください。',
        prompt: {
          model: '',
          type: 'text'
        },
        cancel: true
      })
        .onOk((name: string) => {
          changeName(name);
        });
    };

    const deleteAccount = (router: Router) => {
      api.delete('/api/users/me')
        .then(() => {
          $q.notify({
            color: 'positive',
            position: 'bottom',
            message: 'アカウントを削除しました。',
            icon: 'done'
          });
          router.push('/login').catch((err) => {
            console.log('router push failed', err)
          })
        })
        .catch(() => {
          $q.notify({
            color: 'negative',
            position: 'bottom',
            message: 'アカウントの削除に失敗しました。',
            icon: 'report_problem'
          });
        });
    };
    const deleteConfirm = (router: Router) => {
      $q.dialog({
        title: 'アカウントを削除する',
        message: '本当にアカウントを削除しますか？',
        cancel: true
      })
        .onOk(() => {
          deleteAccount(router);
        });
    };

    reload();

    return { data, loading, changeNamePrompt, deleteConfirm };
  }
});
</script>
