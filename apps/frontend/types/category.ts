// Mirrors internal/domain.VideoCategory in apps/backend — keep in sync if the Go struct changes.
export interface Category {
  id: string;
  countryCode: string;
  title: string;
  assignable: boolean;
}

export interface CategoriesResponse {
  success: boolean;
  data?: Category[];
  error?: string;
}
