package dev.gitswitch.plugin

import com.intellij.openapi.progress.ProgressManager
import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.popup.JBPopupFactory
import com.intellij.openapi.ui.popup.PopupStep
import com.intellij.openapi.ui.popup.util.BaseListPopupStep
import com.intellij.ui.awt.RelativePoint
import com.intellij.openapi.wm.WindowManager

object SwitchProfilePopup {

    fun show(project: Project, at: RelativePoint) {
        val profiles = GitswitchCli.list(project.basePath)
        if (profiles.isEmpty()) {
            notify(project, "gitswitch: no profiles configured. Run 'gitswitch add' first.")
            return
        }
        val step = object : BaseListPopupStep<Profile>("Switch git identity", profiles) {
            override fun getTextFor(value: Profile): String =
                (if (value.active) "✓ " else "  ") + "${value.nickname} — ${value.userName} <${value.email}>"

            override fun onChosen(selectedValue: Profile, finalChoice: Boolean): PopupStep<*>? {
                doInBackground(project) {
                    GitswitchCli.switchTo(project.basePath, selectedValue.nickname)
                    refreshWidget(project)
                }
                return PopupStep.FINAL_CHOICE
            }
        }
        JBPopupFactory.getInstance().createListPopup(step).show(at)
    }

    private fun doInBackground(project: Project, action: () -> Unit) {
        ProgressManager.getInstance().run(object : com.intellij.openapi.progress.Task.Backgroundable(project, "gitswitch", false) {
            override fun run(indicator: com.intellij.openapi.progress.ProgressIndicator) {
                try {
                    action()
                } catch (e: GitswitchCliException) {
                    notify(project, "gitswitch: ${e.message}")
                }
            }
        })
    }

    private fun refreshWidget(project: Project) {
        WindowManager.getInstance().getStatusBar(project)?.updateWidget("gitswitch.statusBarWidget")
    }

    private fun notify(project: Project, message: String) {
        com.intellij.notification.NotificationGroupManager.getInstance()
            .getNotificationGroup("gitswitch")
            .createNotification(message, com.intellij.notification.NotificationType.WARNING)
            .notify(project)
    }
}
