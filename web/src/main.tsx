import { render } from "@solidjs/web";
import App from "./plandalf/App";
import "./plandalf/base.css";

const root = document.getElementById("root");
if (!root) throw new Error("Missing #root");

render(() => <App />, root);
