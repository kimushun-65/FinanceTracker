import { ReportsContainer } from '../../page-components/reports';
import { ProtectedRoute } from '../../shared/ui/components/ProtectedRoute';

export default function ReportsPage() {
  return (
    <ProtectedRoute>
      <ReportsContainer />
    </ProtectedRoute>
  );
}
