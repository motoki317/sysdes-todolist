<template>
  <q-page padding class='center column'>
    <h4>ログイン</h4>

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
        autocomplete='current-password'
        v-model='password'
        label='パスワード'
        lazy-rules
        :rules="[ val => val && val.length > 0 || 'パスワードを入力してください' ]"
      />

      <div class='q-gutter-md' style='display: flex; justify-content: end'>
        <q-btn label='もしくは登録する' color='secondary' @click='$router.push("/register")' style='margin: 10px'/>
        <q-btn label='ログイン' type='submit' color='primary' style='margin: 10px'/>
      </div>
    </q-form>

    <div>
    </div>
  </q-page>
</template>

<script lang='ts'>
import { defineComponent, ref } from 'vue';
import { api } from 'boot/axios';
import { useQuasar } from 'quasar';
import { Router } from 'vue-router'

export default defineComponent({
  name: 'Login',

  props: {},

  setup() {
    const $q = useQuasar()

    const name = ref<string>('');
    const password = ref<string>('');

    const onSubmit = (router: Router) => {
      api.post('/api/login', { name: name.value, password: password.value })
        .then(() => {
          return router.push('/')
        })
        .catch(() => {
          $q.notify({
            color: 'negative',
            position: 'bottom',
            message: 'ログインに失敗しました。',
            icon: 'report_problem'
          })
        })
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
