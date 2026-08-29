import type { JSX } from "solid-js";
import * as stylex from "@stylexjs/stylex";

export type StackProps = {
  children: JSX.Element;
  direction?: "row" | "column";
  gap?: "xs" | "sm" | "md" | "lg";
  align?: "start" | "center" | "stretch";
};

export function Stack(props: StackProps) {
  return (
    <div
      {...stylex.props(
        styles.stack,
        props.direction === "row" && styles.row,
        props.gap === "xs" && styles.gapXs,
        props.gap === "sm" && styles.gapSm,
        (!props.gap || props.gap === "md") && styles.gapMd,
        props.gap === "lg" && styles.gapLg,
        props.align === "start" && styles.alignStart,
        props.align === "center" && styles.alignCenter,
        (!props.align || props.align === "stretch") && styles.alignStretch,
      )}
    >
      {props.children}
    </div>
  );
}

export type TextProps = {
  children: JSX.Element;
  tone?: "default" | "muted" | "accent";
  size?: "sm" | "md" | "lg" | "xl";
  weight?: "regular" | "medium" | "bold";
};

export function Text(props: TextProps) {
  return (
    <span
      {...stylex.props(
        styles.text,
        props.tone === "muted" && styles.muted,
        props.tone === "accent" && styles.accent,
        props.size === "sm" && styles.textSm,
        (!props.size || props.size === "md") && styles.textMd,
        props.size === "lg" && styles.textLg,
        props.size === "xl" && styles.textXl,
        props.weight === "medium" && styles.medium,
        props.weight === "bold" && styles.bold,
      )}
    >
      {props.children}
    </span>
  );
}

export type ButtonProps = {
  children: JSX.Element;
  onClick?: () => void;
  disabled?: boolean;
  tone?: "primary" | "secondary" | "danger";
  wide?: boolean;
};

export function Button(props: ButtonProps) {
  return (
    <button
      type="button"
      disabled={props.disabled}
      onClick={() => props.onClick?.()}
      {...stylex.props(
        styles.button,
        (!props.tone || props.tone === "primary") && styles.buttonPrimary,
        props.tone === "secondary" && styles.buttonSecondary,
        props.tone === "danger" && styles.buttonDanger,
        props.wide && styles.wide,
      )}
    >
      {props.children}
    </button>
  );
}

export function Surface(props: { children: JSX.Element }) {
  return <section {...stylex.props(styles.surface)}>{props.children}</section>;
}

const styles = stylex.create({
  stack: {
    display: "flex",
    flexDirection: "column",
  },
  row: {
    flexDirection: "row",
  },
  gapXs: { gap: 6 },
  gapSm: { gap: 10 },
  gapMd: { gap: 16 },
  gapLg: { gap: 24 },
  alignStart: { alignItems: "flex-start" },
  alignCenter: { alignItems: "center" },
  alignStretch: { alignItems: "stretch" },
  text: {
    color: "#f5f1e8",
    fontFamily: "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, sans-serif",
    lineHeight: 1.45,
  },
  muted: { color: "#a7a397" },
  accent: { color: "#d9bf77" },
  textSm: { fontSize: 13 },
  textMd: { fontSize: 16 },
  textLg: { fontSize: 20 },
  textXl: { fontSize: 30, lineHeight: 1.15 },
  medium: { fontWeight: 600 },
  bold: { fontWeight: 750 },
  button: {
    appearance: "none",
    borderStyle: "solid",
    borderWidth: 1,
    borderRadius: 14,
    cursor: {
      default: "pointer",
      ":disabled": "not-allowed",
    },
    fontFamily: "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, sans-serif",
    fontSize: 15,
    fontWeight: 700,
    minHeight: 48,
    paddingBlock: 12,
    paddingInline: 16,
    opacity: {
      default: 1,
      ":disabled": 0.5,
    },
    transition: "transform 120ms ease, opacity 120ms ease, background-color 120ms ease",
    transform: {
      default: "translateY(0)",
      ":active": "translateY(1px)",
    },
  },
  buttonPrimary: {
    backgroundColor: {
      default: "#d9bf77",
      ":hover": "#e5cc85",
    },
    borderColor: "#d9bf77",
    color: "#181713",
  },
  buttonSecondary: {
    backgroundColor: {
      default: "#24231f",
      ":hover": "#2d2b26",
    },
    borderColor: "#424038",
    color: "#f5f1e8",
  },
  buttonDanger: {
    backgroundColor: {
      default: "#3a211f",
      ":hover": "#4a2825",
    },
    borderColor: "#69423d",
    color: "#ffd8d2",
  },
  wide: { width: "100%" },
  surface: {
    backgroundColor: "#1d1c19",
    borderColor: "#35332d",
    borderRadius: 22,
    borderStyle: "solid",
    borderWidth: 1,
    boxShadow: "0 18px 60px rgba(0, 0, 0, 0.22)",
    padding: 20,
  },
});
