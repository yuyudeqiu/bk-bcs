<template>
  <span v-if="flagsMap.BKAI">
    <!-- AI小鲸按钮 -->
    <bcs-popover theme="ai-assistant light" placement="bottom" trigger="manual" offset="0, 10" ref="popoverRef">
      <span
        class="relative top-[2px] flex items-center justify-center w-[18px] h-[18px] text-[14px] cursor-pointer"
        v-bk-tooltips="$t('blueking.aiScriptsAssistant.desc')"
        @click="handleShowAssistant">
        <img :src="AssistantSmallIcon" />
      </span>
      <template #content>
        <div
          :class="[
            'bg-[#fff] px-[16px]',
            'flex items-center h-[40px]  rounded-[20px]',
            'shadow-[0_2px_6px_0_rgba(0,0,0,0.16)] hover:shadow-[0_2px_8px_0_rgba(0,0,0,0.2)]'
          ]">
          <img :src="AssistantIcon" />
          {{ $t('blueking.aiScriptsAssistant.errTips') }}
        </div>
      </template>
    </bcs-popover>
    <!-- AI小鲸对话框 -->
    <Assistant
      :is-show.sync="isShowAssistant"
      :loading="loading"
      :messages="messages"
      :position-limit="positionLimit"
      :prompts="prompts"
      :start-position="startPosition"
      :size-limit="sizeLimit"
      enable-popup
      @clear="handleClear"
      @close="handleClose"
      @send="handleSend"
      @stop="handleStop" />
  </span>
</template>
<script setup lang="ts">
import { debounce } from 'lodash';
import { ref } from 'vue';

import Assistant, { ChatHelper, IMessage, ISendData, IStartPosition, MessageStatus, RoleType } from '@blueking/ai-blueking/vue2';

import '@blueking/ai-blueking/dist/vue2/style.css';
import { BCS_UI_PREFIX } from '@/common/constant';
import { useAppData } from '@/composables/use-app';
import $i18n from '@/i18n/i18n-setup';
import AssistantIcon from '@/images/assistant.png';
import AssistantSmallIcon from '@/images/assistant-small.svg';

// 属性配置
const props = defineProps({
  preset: {
    type: String,
    default: '',
  },
});

// 特性开关
const { flagsMap } = useAppData();

// AI小鲸
const streamID = ref(1);
const popoverRef = ref();
const isShowAssistant = ref(false);
const loading = ref(false);
const messages = ref<IMessage[]>([]);
const prompts = ref([]);
const positionLimit = ref({
  top: 0,
  bottom: 0,
  left: 0,
  right: 0,
});
const sizeLimit = ref({
  height: 460,
  width: 720,
});
const startPosition = ref<IStartPosition>({
  right: 24,
  top: window.innerHeight - sizeLimit.value.height - 10,
  bottom: 10,
  left: window.innerWidth - sizeLimit.value.width - 24,
});
// 清空消息
const handleClear = () => {
  messages.value.splice(0);
};

// 聊天开始
const handleStart = () => {
  loading.value = true;
  messages.value.push({
    role: RoleType.Assistant,
    content: $i18n.t('blueking.aiScriptsAssistant.loading'),
    status: MessageStatus.Loading,
  });
};

// 接收消息
const handleReceiveMessage = (msg: string, id: number | string, cover?: boolean) => {
  const currentMessage = messages.value.at(-1);
  if (!currentMessage) return;

  if (currentMessage.status === 'loading') {
    // 如果是loading状态，直接覆盖
    currentMessage.content = msg;
    currentMessage.status = MessageStatus.Success;
  } else if (currentMessage.status === 'success') {
    // 如果是后续消息，就追加消息
    currentMessage.content = cover ? msg : currentMessage.content + msg;
  }
};

// 聊天结束
const handleEnd = () => {
  loading.value = false;
  const currentMessage = messages.value.at(-1);
  if (!currentMessage) return;
  // loading 情况下终止
  if (currentMessage.status === 'loading') {
    currentMessage.content = $i18n.t('blueking.aiScriptsAssistant.breakLoading');
    currentMessage.status = MessageStatus.Error;
  }
};

// 终止聊天
const handleStop = async () => {
  await chatHelper.stop(streamID.value);
};

// 错误处理
const handleError = (msg: string) => {
  const currentMessage = messages.value.at(-1);
  if (!currentMessage) return;

  currentMessage.status = MessageStatus.Error;
  currentMessage.content = msg;
  loading.value = false;
};
const chatHelper = new ChatHelper(`${BCS_UI_PREFIX}/assistant`, handleStart, handleReceiveMessage, handleEnd, handleError, messages.value);
// 发送消息
const handleSend = async (args: ISendData) => {
  if (!flagsMap.value.BKAI || !props.preset) return;

  // 记录当前消息记录
  const chatHistory = [...messages.value];

  // 添加用户消息
  messages.value.push({
    role: RoleType.User,
    content: args.content,
    cite: args.cite,
  });
  // 根据参数构造输入内容
  // eslint-disable-next-line no-nested-ternary
  const input = args.prompt
    ? args.prompt                           // 如果有 prompt，直接使用
    : args.cite
      ? `${args.content}: ${args.cite}`     // 如果有 cite，拼接 content 和 cite
      : args.content;                       // 否则只使用 content
  // ai 消息，id是唯一标识当前流，调用 chatHelper.stop 的时候需要传入
  chatHelper.stream({
    inputs: {
      input,
      chat_history: chatHistory,
      preset: props.preset,
    },
  }, streamID.value);
};
// 发送消息防抖(外部调用)
const handleSendMsg = debounce((msg: string) => {
  handleSend({ content: msg });
}, 1000);

// 关闭对话框
const handleClose = () => {
  isShowAssistant.value = false;
};
// 快捷prompt(暂时不启用改功能)
// const handleChoosePrompt = (prompt) => {
//   console.log(prompt);
// };
// 显示对话框
const handleShowAssistant = () => {
  isShowAssistant.value = true;
};
// 消息提示
const showAITips = () => {
  if (!flagsMap.value.BKAI) return;
  popoverRef.value?.showHandler();
};

defineExpose({
  handleSendMsg,
  showAITips,
});
</script>
<style lang="postcss">
.tippy-tooltip.ai-assistant-theme {
  padding: 0!important;
  box-shadow: unset !important;
  background: transparent;
}
</style>
