type PageHeaderProps = {
  title: string;
  subtitle: string;
};

export function PageHeader({ title, subtitle }: PageHeaderProps) {
  return (
    <div className="p-4 pl-7">
      <h1 className="text-3xl font-bold">{title}</h1>
      <span className="text-sm text-gray-400">{subtitle}</span>
    </div>
  );
}
