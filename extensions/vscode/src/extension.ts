import * as vscode from "vscode";
import { execFile } from "child_process";
import { promisify } from "util";

const execFileAsync = promisify(execFile);

interface CurrentProfile {
  nickname?: string;
  user_name?: string;
  email?: string;
  scope: string;
}

interface Profile {
  nickname: string;
  user_name: string;
  email: string;
  active: boolean;
}

function bin(): string {
  return vscode.workspace.getConfiguration("gitswitch").get<string>("binaryPath") || "gitswitch";
}

function cwd(): string | undefined {
  return vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
}

async function run(args: string[]): Promise<string> {
  const { stdout } = await execFileAsync(bin(), args, { cwd: cwd() });
  return stdout;
}

async function getCurrent(): Promise<CurrentProfile | undefined> {
  try {
    return JSON.parse(await run(["current", "--json"]));
  } catch {
    return undefined;
  }
}

async function getProfiles(): Promise<Profile[]> {
  try {
    return JSON.parse(await run(["list", "--json"]));
  } catch {
    return [];
  }
}

const SCOPE_LABEL: Record<string, string> = { repo: "pinned", session: "session", global: "global" };

export function activate(context: vscode.ExtensionContext) {
  const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);
  statusBarItem.command = "gitswitch.switchProfile";
  context.subscriptions.push(statusBarItem);

  async function refresh() {
    const current = await getCurrent();
    if (!current) {
      statusBarItem.hide();
      return;
    }
    if (!current.nickname) {
      statusBarItem.text = "$(person) gitswitch: no profile";
      statusBarItem.tooltip = "Click to switch git identity";
    } else {
      statusBarItem.text = `$(person) ${current.nickname}`;
      statusBarItem.tooltip = `${current.user_name} <${current.email}> (${SCOPE_LABEL[current.scope] ?? current.scope})\nClick to switch git identity`;
    }
    statusBarItem.show();
  }

  context.subscriptions.push(
    vscode.commands.registerCommand("gitswitch.switchProfile", async () => {
      const profiles = await getProfiles();
      if (profiles.length === 0) {
        vscode.window.showWarningMessage("gitswitch: no profiles configured. Run 'gitswitch add' first.");
        return;
      }
      const pick = await vscode.window.showQuickPick(
        profiles.map((p) => ({
          label: p.active ? `$(check) ${p.nickname}` : p.nickname,
          description: `${p.user_name} <${p.email}>`,
          nickname: p.nickname,
        })),
        { placeHolder: "Select a git identity to switch to" }
      );
      if (!pick) return;
      try {
        await run(["switch", pick.nickname]);
        await refresh();
      } catch (err) {
        vscode.window.showErrorMessage(`gitswitch switch failed: ${(err as Error).message}`);
      }
    }),

    vscode.commands.registerCommand("gitswitch.pinProfile", async () => {
      if (!cwd()) {
        vscode.window.showWarningMessage("gitswitch: open a folder first.");
        return;
      }
      const profiles = await getProfiles();
      const pick = await vscode.window.showQuickPick(
        profiles.map((p) => ({ label: p.nickname, description: `${p.user_name} <${p.email}>`, nickname: p.nickname })),
        { placeHolder: "Pin which identity to this repo?" }
      );
      if (!pick) return;
      try {
        await run(["pin", pick.nickname]);
        await refresh();
        vscode.window.showInformationMessage(`gitswitch: pinned '${pick.nickname}' to this repo`);
      } catch (err) {
        vscode.window.showErrorMessage(`gitswitch pin failed: ${(err as Error).message}`);
      }
    }),

    vscode.commands.registerCommand("gitswitch.unpin", async () => {
      try {
        await run(["unpin"]);
        await refresh();
        vscode.window.showInformationMessage("gitswitch: unpinned this repo");
      } catch (err) {
        vscode.window.showErrorMessage(`gitswitch unpin failed: ${(err as Error).message}`);
      }
    })
  );

  refresh();
  // Catches identity changes made outside VS Code (e.g. `gitswitch switch` in a terminal).
  context.subscriptions.push(vscode.window.onDidChangeWindowState((s) => s.focused && refresh()));
  context.subscriptions.push(vscode.workspace.onDidChangeWorkspaceFolders(refresh));
}

export function deactivate() {}
