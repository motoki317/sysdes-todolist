export interface User {
  name: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface Todo {
  id: number;
  title: string;
  done: boolean;
  createdAt: Date;
  updatedAt: Date;
}
