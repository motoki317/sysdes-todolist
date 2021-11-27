<template>
  <q-page padding class='center column'>
    <h4>新規登録</h4>

    <q-form
      @submit='onSubmit($router)'
      class='q-gutter-md'
      style='width: min(80vw, 500px)'
    >
      <q-input
        filled
        type='text'
        autocomplete='username'
        v-model='name'
        label='ユーザー名'
        lazy-rules
        :rules="[ val => val && val.length > 0 || 'ユーザー名を入力してください']"
      />

      <q-input
        filled
        type='password'
        autocomplete='new-password'
        v-model='password'
        label='パスワード'
        lazy-rules
        :rules="[ val => val && val.length > 0 || 'パスワードを入力してください' ]"
      />

      <div class='q-gutter-md' style='display: flex; justify-content: end'>
        <q-btn label='ログイン画面へ戻る' color='secondary' @click='$router.push("/login")' style='margin: 10px'/>
        <q-btn label='登録' type='submit' color='primary' style='margin: 10px' />
      </div>
    </q-form>

    <div>
    </div>
  </q-page>
</template>

<script lang='ts'>
import { defineComponent, ref } from 'vue';
import { useQuasar } from 'quasar';
import { Router } from 'vue-router';
import { api } from 'boot/axios';
import { AxiosError } from 'axios';

export default defineComponent({
  name: 'Register',

  props: {},

  setup() {
    const $q = useQuasar();

    const name = ref<string>('');
    const password = ref<string>('');

    const onSubmit = (router: Router) => {
      api.post('/api/signup', { name: name.value, password: password.value })
        .then(() => {
          $q.notify({
            color: 'positive',
            position: 'bottom',
            message: '登録に成功しました。ログインしてください。',
            icon: 'done'
          });
          return router.push('/login');
        })
        .catch((e) => {
          if ((e as AxiosError).response?.status === 409) {
            $q.notify({
              color: 'negative',
              position: 'bottom',
              message: 'すでに使われているユーザー名です。',
              icon: 'report_problem'
            })
          } else {
            $q.notify({
              color: 'negative',
              position: 'bottom',
              message: '登録に失敗しました。',
              icon: 'report_problem'
            });
          }
        });
    };

    return { name, password, onSubmit };
  }
});
</script>

<style>
.center {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
