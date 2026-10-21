package dev.gitswitch.plugin

import com.google.gson.Gson
import com.google.gson.annotations.SerializedName
import com.intellij.openapi.diagnostic.logger
import java.io.File

data class CurrentProfile(
    val nickname: String? = null,
    @SerializedName("user_name") val userName: String? = null,
    val email: String? = null,
    val scope: String = "global",
    @SerializedName("guard_active") val guardActive: Boolean = false,
)

data class Profile(
    val nickname: String,
    @SerializedName("user_name") val userName: String,
    val email: String,
    val active: Boolean,
)

class GitswitchCliException(message: String) : Exception(message)

object GitswitchCli {
    private val log = logger<GitswitchCli>()
    private val gson = Gson()

    private fun binaryPath(): String =
        System.getProperty("gitswitch.binaryPath") ?: "gitswitch"

    private fun run(workDir: String?, vararg args: String): String {
        val cmd = listOf(binaryPath()) + args
        val builder = ProcessBuilder(cmd)
        if (workDir != null) builder.directory(File(workDir))
        val process = builder.start()
        val stdout = process.inputStream.bufferedReader().readText()
        val stderr = process.errorStream.bufferedReader().readText()
        val exit = process.waitFor()
        if (exit != 0) {
            log.warn("gitswitch ${args.joinToString(" ")} failed: $stderr")
            throw GitswitchCliException(stderr.ifBlank { "gitswitch exited with code $exit" })
        }
        return stdout
    }

    fun current(workDir: String?): CurrentProfile? = try {
        gson.fromJson(run(workDir, "current", "--json"), CurrentProfile::class.java)
    } catch (e: Exception) {
        log.debug("gitswitch current failed", e)
        null
    }

    fun list(workDir: String?): List<Profile> = try {
        gson.fromJson(run(workDir, "list", "--json"), Array<Profile>::class.java).toList()
    } catch (e: Exception) {
        log.debug("gitswitch list failed", e)
        emptyList()
    }

    fun switchTo(workDir: String?, nickname: String) = run(workDir, "switch", nickname)

    fun pin(workDir: String?, nickname: String) = run(workDir, "pin", nickname)

    fun unpin(workDir: String?) = run(workDir, "unpin")

    fun setGuard(workDir: String?, enable: Boolean) =
        run(workDir, "guard", if (enable) "enable" else "disable")
}
