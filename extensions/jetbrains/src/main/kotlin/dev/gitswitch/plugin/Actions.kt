package dev.gitswitch.plugin

import com.intellij.notification.NotificationGroupManager
import com.intellij.notification.NotificationType
import com.intellij.openapi.actionSystem.ActionUpdateThread
import com.intellij.openapi.actionSystem.AnAction
import com.intellij.openapi.actionSystem.AnActionEvent
import com.intellij.openapi.progress.ProgressIndicator
import com.intellij.openapi.progress.ProgressManager
import com.intellij.openapi.progress.Task
import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.popup.JBPopupFactory
import com.intellij.openapi.ui.popup.PopupStep
import com.intellij.openapi.ui.popup.util.BaseListPopupStep
import com.intellij.openapi.wm.WindowManager

private fun notify(project: Project, message: String, type: NotificationType = NotificationType.INFORMATION) {
    NotificationGroupManager.getInstance()
        .getNotificationGroup("gitswitch")
        .createNotification(message, type)
        .notify(project)
}

private fun refreshWidget(project: Project) {
    WindowManager.getInstance().getStatusBar(project)?.updateWidget("gitswitch.statusBarWidget")
}

private fun runInBackground(project: Project, action: () -> Unit) {
    ProgressManager.getInstance().run(object : Task.Backgroundable(project, "gitswitch", false) {
        override fun run(indicator: ProgressIndicator) {
            try {
                action()
            } catch (e: GitswitchCliException) {
                notify(project, "gitswitch: ${e.message}", NotificationType.WARNING)
            }
        }
    })
}

class SwitchProfileAction : AnAction("Switch Profile") {
    override fun getActionUpdateThread() = ActionUpdateThread.BGT
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.project ?: return
        val statusBar = WindowManager.getInstance().getStatusBar(project) ?: return
        val anchor = statusBar.component ?: return
        SwitchProfilePopup.show(project, com.intellij.ui.awt.RelativePoint.getCenterOf(anchor))
    }
}

class PinProfileAction : AnAction("Pin Profile to This Repo") {
    override fun getActionUpdateThread() = ActionUpdateThread.BGT
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.project ?: return
        val workDir = project.basePath
        val profiles = GitswitchCli.list(workDir)
        if (profiles.isEmpty()) {
            notify(project, "gitswitch: no profiles configured. Run 'gitswitch add' first.", NotificationType.WARNING)
            return
        }
        val step = object : BaseListPopupStep<Profile>("Pin which identity to this repo?", profiles) {
            override fun getTextFor(value: Profile): String = "${value.nickname} — ${value.userName} <${value.email}>"
            override fun onChosen(selectedValue: Profile, finalChoice: Boolean): PopupStep<*>? {
                runInBackground(project) {
                    GitswitchCli.pin(workDir, selectedValue.nickname)
                    refreshWidget(project)
                    notify(project, "gitswitch: pinned '${selectedValue.nickname}' to this repo")
                }
                return PopupStep.FINAL_CHOICE
            }
        }
        JBPopupFactory.getInstance().createListPopup(step).showCenteredInCurrentWindow(project)
    }
}

class UnpinAction : AnAction("Unpin This Repo") {
    override fun getActionUpdateThread() = ActionUpdateThread.BGT
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.project ?: return
        val workDir = project.basePath
        runInBackground(project) {
            GitswitchCli.unpin(workDir)
            refreshWidget(project)
            notify(project, "gitswitch: unpinned this repo")
        }
    }
}

class ToggleGuardAction : AnAction("Toggle Pre-commit Guard") {
    override fun getActionUpdateThread() = ActionUpdateThread.BGT
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.project ?: return
        val workDir = project.basePath
        runInBackground(project) {
            val current = GitswitchCli.current(workDir)
            val enable = !(current?.guardActive ?: false)
            GitswitchCli.setGuard(workDir, enable)
            refreshWidget(project)
            notify(project, "gitswitch: pre-commit guard ${if (enable) "enabled" else "disabled"}")
        }
    }
}
