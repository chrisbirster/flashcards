type DeckStatItemProps = {
    stat: number;
    label: string;
    color: string;
};

export function DeckStatItem({ stat, label, color }: DeckStatItemProps) {
    const mergedClassName = `${color} font-medium`;
    return (
        <span className={mergedClassName}>
            {stat} {label}
        </span>
    )
}