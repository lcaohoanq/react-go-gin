import { Flex, Text, IconButton, Badge, Spinner } from "@chakra-ui/react";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Todo } from "./TodoList";
import { todoApi } from "../api/axios";

const TodoItem = ({ todo }: { todo: Todo }) => {
  const queryClient = useQueryClient();
  const { mutate: updateTodo, isPending: isUpdating } = useMutation({
    mutationFn: async () => {
      const res = await todoApi.updateTodo(todo.ID);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });

  const { mutate: deleteTodo, isPending: isDeleting } = useMutation({
    mutationFn: async () => {
      const res = await todoApi.deleteTodo(todo.ID);
      if (res.status !== 204) {
        throw new Error("Failed to delete todo");
      }
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });

  return (
    <Flex
      direction="column"
      p={4}
      shadow="md"
      borderWidth="1px"
      borderRadius="lg"
      justify="space-between"
      align="center"
      className="w-fit"
    >
      <Text
        color={todo.completed ? "gray.500" : "inherit"}
        textDecoration={todo.completed ? "line-through" : "none"}
      >
        {todo.body}
      </Text>
      <Flex gap={2}>
        <Badge
          colorScheme={todo.completed ? "green" : "yellow"}
          p={2}
          borderRadius="md"
        >
          {todo.completed ? "Completed" : "Pending"}
        </Badge>
        <IconButton
          aria-label="Mark as complete"
          icon={isUpdating ? <Spinner size="sm" /> : <FaCheckCircle />}
          colorScheme="green"
          variant="ghost"
          isDisabled={todo.completed}
          onClick={() => updateTodo()}
        />
        <IconButton
          aria-label="Delete todo"
          icon={isDeleting ? <Spinner size="sm" /> : <MdDelete />}
          colorScheme="red"
          variant="ghost"
          onClick={() => deleteTodo()}
        />
      </Flex>
    </Flex>
  );
};

export default TodoItem;
