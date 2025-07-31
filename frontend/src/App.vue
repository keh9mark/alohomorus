<template>
  <div class="main">
    <div class="main__top"></div>
    <div class="main__center">
      <div class="main__left">
        <div class="left__top">
          <div class="_right">
            <img src="./assets/images/energy.png" />
            <div class="_label">EnergyCalc</div>
          </div>
          <div class="_version">v0.0.1</div>
        </div>
        <div class="_menu">
          <div class="left__menu">
            <div class="menu_item disabled">
              <img src="./assets/images/list.png" />
              <div class="label">Список расчетов</div>
            </div>
            <div class="menu_item clicked">
              <img src="./assets/images/add-item.png" />
              <div class="label">Редактор</div>
            </div>
            <div class="menu_item disabled">
              <img src="./assets/images/download.png" />
              <div class="label">Выгрузить данные</div>
            </div>
            <div class="menu_item disabled">
              <img src="./assets/images/upload.png" />
              <div class="label">Загрузить данные</div>
            </div>
          </div>
          <div class="left_bottom_menu">
            <div class="menu_item disabled">
              <img src="./assets/images/help.png" />
              <div class="label">Информация</div>
            </div>
            <div class="menu_item disabled">
              <img src="./assets/images/settings.png" />
              <div class="label">Настройки</div>
            </div>
          </div>
        </div>
      </div>
      <div class="main__content">
        <editor />
      </div>
      <div class="main__right"></div>
    </div>
    <div class="main__bottom">EnergyCalc @ 2024 Тест-проект v0.0.1</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import NewTask from './components/NewTask.vue'
import editor from "./components/editor/main.vue"

export default defineComponent({
  components: {
    NewTask: NewTask,
    editor: editor
  },
  data() {
    return {}
  },
})
</script>

<style lang="scss">
.main {
  cursor: default;
  height: 100%;
  width: 100%;
  background-color: #fafafa;
  color: black;

  >.main__top {
    height: 0px;
    width: 100%;
  }

  >.main__center {
    height: calc(100% - 40px);
    width: calc(100% - 10px);
    display: flex;

    >.main__left {
      padding-top: 20px;
      height: 100%;
      width: 300px;

      >.left__top {
        height: 50px;

        display: flex;
        align-items: center;
        justify-content: space-between;

        >._right {
          padding-left: 20px;
          display: flex;
          align-items: center;

          >img {
            height: 50px;
            width: 50px;
          }

          >._label {
            padding-left: 10px;
          }
        }

        >._version {
          padding-right: 20px;
        }
      }

      >._menu {
        display: flex;
        flex-direction: column;
        height: calc(100% - 100px);
        justify-content: space-between;

        >.left__menu {
          padding-top: 30px;
        }

        >.left_bottom_menu {
          height: 100px;
        }
      }
    }

    >.main__content {
      background-color: white;
      height: calc(100% - 40px);
      width: calc(100% - 350px);
      margin: 20px;
      border: 1px solid #eae4e4;
      border-radius: 20px;
    }
  }

  >.main__bottom {
    height: 30px;
    width: 100%;
  }

  .menu_item {
    padding-left: 30px;
    padding-top: 10px;
    padding-bottom: 10px;
    align-items: center;
    display: flex;
    cursor: pointer;

    &.disabled {
      cursor: not-allowed;
      color: #66737369;
    }

    &.clicked {
      background-color: #e4e7e7;
    }

    >img {
      height: 30px;
      width: 30px;
    }

    >.label {
      padding-left: 10px;
    }
  }
}
</style>
