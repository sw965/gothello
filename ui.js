const boardContainer = document.getElementById("board-container");

const BOARD_ROW = 8;
const BOARD_COLUMN = 8;

const COLUMNS = ["a", "b", "c", "d", "e", "f", "g", "h"];

// ボードを初期化する関数
function initializeBoard() {

    // 上の列番号を追加
    const emptyCell = document.createElement("div");
    emptyCell.classList.add("header-cell");
    boardContainer.appendChild(emptyCell);

    COLUMNS.forEach(col => {
        const headerCell = document.createElement("div");
        headerCell.textContent = col;
        headerCell.classList.add("header-cell");
        boardContainer.appendChild(headerCell);
    });

    // 行番号とセルを作成
    for (let i = 0; i < BOARD_ROW; i++) {
        const rowHeader = document.createElement("div");
        rowHeader.textContent = i + 1;
        rowHeader.classList.add("header-cell");
        boardContainer.appendChild(rowHeader);

        for (let j = 0; j < BOARD_COLUMN; j++) {
            const cellContainer = document.createElement("div");
            cellContainer.classList.add("cell");
            cellContainer.dataset.row = i;
            cellContainer.dataset.col = j;
            boardContainer.appendChild(cellContainer);
        }
    }
}

async function fetchCurrentTurn() {
    const url = new URL(QUERY_URL);
    url.searchParams.append("query_type", encodeURIComponent("current_turn"));
    const response = await fetch(url.toString());
    const currentTurn = await response.json();
    return currentTurn;
}

async function fetchLegalActions() {
    const url = new URL(QUERY_URL);
    url.searchParams.append("query_type", encodeURIComponent("legal_actions"));
    const response = await fetch(url.toString());
    const legalActions = await response.json();
    return Array.from(legalActions);
}

async function fetchIsEnd() {
    const url = new URL(QUERY_URL);
    url.searchParams.append("query_type", encodeURIComponent("is_end"));
    const response = await fetch(url.toString());
    const isEnd = await response.json();
    return isEnd;
}

async function fetchBlackCount() {
    const url = new URL(QUERY_URL);
    url.searchParams.append("query_type", encodeURIComponent("black_count"));
    const response = await fetch(url.toString());
    const winner = response.json();
    return winner;
}

async function fetchWhiteCount() {
    const url = new URL(QUERY_URL);
    url.searchParams.append("query_type", encodeURIComponent("white_count"));
    const response = await fetch(url.toString());
    const count = response.json();
    return count;
}

async function fetchMCTSResult() {
    const simulations = document.getElementById("simulations");
    const url = new URL(QUERY_URL);
    url.searchParams.append("query_type", encodeURIComponent("mcts_result"));
    url.searchParams.append("simulation", encodeURIComponent(simulations.value));
    const response = await fetch(url.toString())
    const mctsResult = response.json();
    return mctsResult;
}

async function fetchState() {
    const url = new URL(QUERY_URL);
    url.searchParams.append("query_type", encodeURIComponent("state"));
    const response = await fetch(url.toString())
    const json = await response.json();
    return JSON.parse(JSON.stringify(json));
}

async function push(action) {
    const url = new URL(COMMAND_URL);
    url.searchParams.append("command_type", encodeURIComponent("push"));
    url.searchParams.append("action", encodeURIComponent(action.Row + " " + action.Column));
    await fetch(url.toString());
    return;
}

// 全てのセルをクリック不可にする関数
function disableAllCells() {
    const cells = document.querySelectorAll(".cell");
    cells.forEach(cell => {
        cell.removeEventListener("click", handleCellClick);
        cell.classList.remove("legal-action");
    });
}

// 合法着手点を有効化する関数
function enableLegalActions(legalActions, currentTurnColor) {
    legalActions.forEach(action => {
        const row = action.Row;
        const col = action.Column;
        const cell = document.querySelector(`.cell[data-row="${row}"][data-col="${col}"]`);
        cell.classList.add("legal-action", currentTurnColor);
        cell.addEventListener("click", handleCellClick);
    });
}

// セルがクリックされたときの処理
async function handleCellClick(event) {
    const cell = event.target;

    // 全てのセルを押せなくする
    disableAllCells();

    // 石を置く (仮に黒石を置く)
    cell.classList.add(humanColor);

    await push({Row:cell.dataset.row, Column:cell.dataset.col});
    const state = await fetchState();
    const legalActions = await fetchLegalActions();
    updateBoard(state.Board, legalActions);

    let continuous = 0;
    //人間の手番が回ってくるまでAIが指す。
    while (true) {
        const isEnd = await fetchIsEnd();
        if (isEnd) {
            const blackCount = await fetchBlackCount();
            const whiteCount = await fetchWhiteCount();
            const msg = "黒石：" + blackCount + "\n" + "白石：" + whiteCount
            let humanMsg;
            console.log(blackCount, whiteCount, msg, humanColor);
            if (blackCount > whiteCount) {
                if (humanColor === "black") {
                    humanMsg = "あなたの勝ちです！"
                } else {
                    humanMsg = "AIの勝ちです!"
                }
            } else if (blackCount < whiteCount) {
                if (humanColor === "black") {
                    humanMsg = "AIの勝ちです！"
                } else {
                    humanMsg = "あなたの勝ちです！"
                }
            }
            alert(msg + "\n" + humanMsg);
            const startButton = document.getElementById("start-button");
            startButton.style.display = "block";
            break
        }

        const currentTurn = await fetchCurrentTurn();
        //currentTurnは数値型で、1(黒)か2(白)が格納されている。
        const turnColor = ["black", "white"][currentTurn-1]
        if (turnColor === humanColor) {
            break
        }
        
        if (continuous >= 1) {
            alert("打てる所がないため、パスになります。");
        };

        const mctsResult = await fetchMCTSResult();
        await push(mctsResult.Action);
        const state = await fetchState();
        const legalActions = await fetchLegalActions();
        updateBoard(Array.from(state.Board), legalActions);

        let humanStoneNumber;
        let aiStoneNumber;
        if (humanColor === "black") {
            humanStoneNumber = 1
            aiStoneNumber = 2
        } else {
            humanStoneNumber = 2
            aiStoneNumber = 1
        }
        updateEvaluationBar(mctsResult.AgentEvals[humanStoneNumber], mctsResult.AgentEvals[aiStoneNumber]);
        continuous += 1
    }
}

// 配列データを受け取ってUIを更新する関数
async function updateBoard(board, legalActions) {
    disableAllCells();
    const cells = document.querySelectorAll(".cell");

    // 全てのセルをリセット
    cells.forEach(cell => {
        cell.classList.remove("black", "white", "legal-action", "black", "white");
        const row = parseInt(cell.dataset.row, 10);
        const col = parseInt(cell.dataset.col, 10);
        if (board[row][col] === 1) {
            cell.classList.add("black");
        } else if (board[row][col] === 2) {
            cell.classList.add("white");
        }
    });

    if (legalActions === undefined) {
        return;
    }

    // 現在の手番を取得
    const currentTurn = await fetchCurrentTurn();
    const turnColor = ["black", "white"][currentTurn - 1];

    // 人間の手番の場合のみ合法手を有効化
    if (turnColor === humanColor) {
        enableLegalActions(legalActions, turnColor);
    }
}

// 評価値バーを更新する関数
function updateEvaluationBar(humanScore, aiScore) {
    const total = humanScore + aiScore;

    const humanPercentage = (humanScore / total) * 100;
    const aiPercentage = (aiScore / total) * 100;

    document.getElementById('human-bar').style.width = humanPercentage + '%';
    document.getElementById('ai-bar').style.width = aiPercentage + '%';

    document.getElementById('human-score').textContent = `人間: ${(humanScore*100.0).toFixed(1)}`;
    document.getElementById('ai-score').textContent = `AI: ${(aiScore*100.0).toFixed(1)}`;
}

const COMMAND_URL = "http://localhost:8064/gothello_command/"
const QUERY_URL = "http://localhost:8064/gothello_query/"
let humanColor;

document.addEventListener("DOMContentLoaded", async () => {
    initializeBoard();
    const startButton = document.getElementById("start-button");

    async function startGame(event) {
        startButton.style.display = "none";

        updateBoard(INIT_BOARD);

        const handRadio = document.querySelector('input[name="hand"]:checked');
        if (handRadio.value === "あなたが先手") {
            humanColor = "black";
        } else {
            humanColor = "white";
        }

        const url = new URL(COMMAND_URL);
        url.searchParams.append("command_type", encodeURIComponent("init"));
        url.searchParams.append("human_color", encodeURIComponent(humanColor));

        fetch(url.toString())
            .then(response => {
                return response.json();
            })
            .then(json => {
                console.log(json);
            })
            .catch(err => {
                console.error('エラー:', err);
            });

        if (humanColor === "white") {
            const mctsResult = await fetchMCTSResult();
            console.log(mctsResult, mctsResult.Action, mctsResult.AgentEvals);
            await push(mctsResult.Action);
            const state = await fetchState();
            const legalActions = await fetchLegalActions();
            updateBoard(state.Board, legalActions);
        } else {
            const legalActions = await fetchLegalActions();
            enableLegalActions(legalActions);
        }
    }
    startButton.addEventListener("click", startGame);

    const INIT_BOARD = [
        [0, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 2, 1, 0, 0, 0],
        [0, 0, 0, 1, 2, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0],
    ];

    updateBoard(INIT_BOARD);
});