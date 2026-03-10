<template>
  <div v-if="show" :class="[styles.modalOverlay, { [styles.closing]: isClosing }]" @animationend="handleAnimationEnd">
    <div :class="styles.modal" @click.stop>
      <div :class="styles.modalHeader">
        <h3>Windows Environment Variables</h3>
        <p :class="styles.modalSubtitle">Understanding how variables and the PATH system are put together</p>

        <div :class="styles.tabNav">
          <button :class="[styles.tabBtn, { [styles.active]: activeTab === 'vars' }]" @click="activeTab = 'vars'">
            Variables
          </button>
          <button :class="[styles.tabBtn, { [styles.active]: activeTab === 'path' }]" @click="activeTab = 'path'">
            PATH
          </button>
        </div>
      </div>

      <div :class="styles.infoContent">
        <!-- Variables Tab -->
        <div v-if="activeTab === 'vars'" :class="styles.tabContent">
          <p :class="styles.introText">Windows uses two distinct types of environment variables to store configuration
            data.</p>

          <div :class="styles.infoSection">
            <h4>User variables</h4>
            <ul>
              <li>Apply to the currently logged-in Windows account</li>
              <li>Only programs run under that account can access them</li>
              <li>No admin rights required for modification</li>
            </ul>
            <span :class="styles.exampleLabel">Example:</span>
            <div :class="styles.exampleBox">
              <code>ANDROID_HOME</code>
              <code>PATH additions for development tools</code>
            </div>
          </div>

          <div :class="styles.divider"></div>

          <div :class="styles.infoSection">
            <h4>System variables</h4>
            <ul>
              <li>Apply to all accounts on the computer</li>
              <li>Used by services, background processes, and all users</li>
              <li>Admin rights are required for modification</li>
            </ul>
            <span :class="styles.exampleLabel">Example:</span>
            <div :class="styles.exampleBox">
              <code>JAVA_HOME for all users</code>
              <code>System-wide utilities</code>
            </div>
          </div>

          <div :class="styles.divider"></div>

          <div :class="styles.infoSection">
            <h4>System Behavior</h4>
            <p>Windows combines both types when a terminal is opened:</p>
            <div :class="styles.logicBox">
              <code>Final PATH = System PATH + User PATH</code>
            </div>
            <p>User PATH entries are appended after the system entries.</p>
          </div>

          <div :class="styles.divider"></div>

          <div :class="styles.infoSection">
            <h4>Recommended Approach</h4>
            <p>For development-related settings (Android SDK, Node.js, etc.):</p>

            <span :class="styles.recommendationLabel">👉 Use User variables</span>
            <span :class="styles.reasonsLabel">Reasons:</span>
            <ul :class="styles.simpleList">
              <li>No admin privileges needed</li>
              <li>Safer (avoids affecting the entire computer)</li>
            </ul>
          </div>

          <div :class="styles.ruleOfThumb">
            <span :class="styles.ruleLabel">Rule of thumb</span>
            <div :class="styles.ruleContent">
              <p>Basically all cases → <strong>User variables</strong></p>
              <p>If a guide specifically mentions otherwise → <strong>System variables</strong></p>
            </div>
          </div>
        </div>

        <!-- PATH Tab -->
        <div v-else :class="styles.tabContent">
          <p :class="styles.introText">The special PATH environment variable specifies the directories Windows searches
            to locate executable programs when a command is entered in a command-line interface.</p>

          <div :class="styles.infoSection">
            <h4>Variable scope</h4>
            <p>The PATH variable can be defined in both <strong>User variables</strong> (affecting only the current
              account) and <strong>System variables</strong> (affecting all accounts).</p>
          </div>

          <div :class="styles.divider"></div>

          <div :class="styles.infoSection">
            <h4>How search works</h4>
            <p>When a command like <code>adb</code> is entered, Windows checks each folder in PATH sequentially. If the
              executable is found in a folder, it is executed. If none of the folders contain the program, a "not
              recognized" error occurs.</p>
          </div>

          <div :class="styles.divider"></div>

          <div :class="styles.infoSection">
            <h4>"Add to PATH"</h4>
            <p>Adding a directory to PATH makes its programs globally available in the terminal without needing to type
              the full file path.</p>

            <span :class="styles.exampleLabel">Instead of typing:</span>
            <div :class="styles.exampleBox">
              <code>C:\...\Android\Sdk\platform-tools\adb.exe</code>
            </div>

            <span :class="styles.exampleLabel">It's enough to type:</span>
            <div :class="styles.exampleBox">
              <code>adb</code>
            </div>
          </div>

          <div :class="styles.divider"></div>

          <div :class="styles.infoSection">
            <h4>Conventions</h4>
            <ul>
              <li><strong>Naming:</strong> The name comes from Unix systems (1970s).</li>
              <li><strong>Casing:</strong> Variables are traditionally uppercase (PATH, ANDROID_HOME, etc.).</li>
            </ul>
          </div>
        </div>
      </div>

      <div :class="styles.modalActions">
        <button :class="[styles.btn, styles.btnPrimary]" @click="$emit('close')">Got it!</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import styles from './EnvInfoModal.module.css'

defineProps({
  show: Boolean,
  isClosing: Boolean
})

const emit = defineEmits(['close', 'animationend'])

const activeTab = ref('vars')

function handleAnimationEnd(e) {
  emit('animationend', e)
}
</script>
