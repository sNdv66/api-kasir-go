export interface User {
  id: string;
  name: string;
  role: "admin" | "kasir";
  branch_id: string;
}