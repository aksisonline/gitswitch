package dev.gitswitch.plugin

import com.intellij.openapi.project.Project
import com.intellij.openapi.wm.StatusBar
import com.intellij.openapi.wm.StatusBarWidget
import com.intellij.openapi.wm.StatusBarWidgetFactory
import com.intellij.openapi.wm.impl.status.EditorBasedWidget
import com.intellij.ui.awt.RelativePoint
import java.awt.event.MouseEvent

private const val WIDGET_ID = "gitswitch.statusBarWidget"

class GitswitchStatusBarWidgetFactory : StatusBarWidgetFactory {
    override fun getId(): String = WIDGET_ID
    override fun getDisplayName(): String = "gitswitch"
    override fun isAvailable(project: Project): Boolean = true
    override fun createWidget(project: Project): StatusBarWidget = GitswitchStatusBarWidget(project)
    override fun disposeWidget(widget: StatusBarWidget) = Unit
    override fun canBeEnabledOn(statusBar: StatusBar): Boolean = true
}

private val SCOPE_LABEL = mapOf("repo" to "pinned", "session" to "session", "global" to "global")

class GitswitchStatusBarWidget(project: Project) : EditorBasedWidget(project), StatusBarWidget.TextPresentation {

    override fun ID(): String = WIDGET_ID
    override fun getPresentation(): StatusBarWidget.WidgetPresentation = this

    override fun getText(): String {
        val current = GitswitchCli.current(project.basePath) ?: return "gitswitch"
        return current.nickname ?: "no profile"
    }

    override fun getTooltipText(): String {
        val current = GitswitchCli.current(project.basePath) ?: return "gitswitch: unavailable"
        if (current.nickname == null) return "gitswitch: no active profile — click to switch"
        val scope = SCOPE_LABEL[current.scope] ?: current.scope
        return "${current.userName} <${current.email}> ($scope)\nClick to switch git identity"
    }

    override fun getAlignment(): Float = java.awt.Component.LEFT_ALIGNMENT

    override fun getClickConsumer() = com.intellij.util.Consumer<MouseEvent> { event ->
        SwitchProfilePopup.show(project, RelativePoint(event))
    }
}
