<template>
  <q-page padding class='column'>

    <h6>パスワード変更</h6>

    <q-form @submit='changePassword($router)' style='width: min(80vw, 500px)'>
      <q-input
        type='password' v-model='currentPassword' autocomplete='current-password' label='現在のパスワード'
        :rules='[val => val && val.length > 0 || "現在のパスワードを入力してください"]'
      ></q-input>
      <q-input
        type='password' v-model='newPassword' autocomplete='new-password' label='新しいパスワード'
        :rules='[val => val && val.length > 0 || "新しいパスワードを入力してください"]'
      ></q-input>
      <q-btn type='submit' label='変更する' color='primary'></q-btn>
    </q-form>

  </q-page>
</template>

<script lang='ts'>
import { defineComponent, ref } from 'vue';
import { api } from 'boot/axios';
import { useQuasar } from 'quasar';
import { AxiosError } from 'axios';
import { Router } from 'vue-router';

export default defineComponent({
  name: 'AccountManagement',

  components: {},

  setup() {
    const $q = useQuasar();

    const currentPassword = ref('');
    const newPassword = ref('');

    const changePassword = (router: Router) => {
      api.put('/api/users/me/password', {
        oldPassword: currentPassword.value,
        newPassword: newPassword.value
      })
        .then(() => {
          $q.notify({
            color: 'positive',
            position: 'bottom',
            message: 'パスワードを変更しました。',
            icon: 'done'
          });
          router.push('/account').catch((err) => {
            console.log('router push failed', err);
          });
        })
        .catch((err) => {
          if ((err as AxiosError).response?.status === 401) {
            $q.notify({
              color: 'negative',
              position: 'bottom',
              message: '現在のパスワードが違います。',
              icon: 'report_problem'
            });
          } else {
            $q.notify({
              color: 'negative',
              position: 'bottom',
              message: 'パスワードの変更に失敗しました。',
              icon: 'report_problem'
            });
          }
        });
    };

    return { currentPassword, newPassword, changePassword };
  }
});
</script>
