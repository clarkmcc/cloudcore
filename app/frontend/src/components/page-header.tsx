import { Skeleton } from "@/components/ui/skeleton.tsx";
import { Button } from "@/components/ui/button.tsx";
import { ArrowLeft } from "lucide-react";
import { useNavigate } from "react-router-dom";

type PageHeaderProps = {
  title: string;
  subtitle: string;
  loading?: boolean;
  backButton?: boolean;
};

export function PageHeader({
  title,
  subtitle,
  loading,
  backButton,
}: PageHeaderProps) {
  const navigate = useNavigate();
  if (loading) {
    return (
      <div className="p-4 pl-7 space-y-3">
        <Skeleton className="h-8 w-1/3" />
        <Skeleton className="h-6 w-1/4" />
      </div>
    );
  }
  return (
    <div className="p-4 pl-7">
      <div className="flex flex-row items-center space-x-2">
        {backButton && (
          <Button variant="secondary" size="icon" onClick={() => navigate(-1)}>
            <ArrowLeft size={15} />
          </Button>
        )}
        <h1 className="text-3xl font-bold">{title}</h1>
      </div>
      <span className="text-sm text-gray-400">{subtitle}</span>
    </div>
  );
}
