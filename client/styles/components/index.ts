export const components = {
  Input: {
    baseStyle: {
      bg: "white",
      color: "text",
      _placeholder: {
        color: "gray.500",
      },
    },
  },
  Button: {
    variants: {
      primary: {
        bg: "primary",
        color: "buttonText",
        _hover: {
          bg: "primary",
          opacity: 0.8,
        },
      },
    },
  },
  VStack: {
    baseStyle: {
      bg: "bg",
      color: "text",
      spacing: 4,
      p: 4,
    },
  },
  HStack: {
    baseStyle: {
      spacing: 2,
    },
  },
  Friend: {
    baseStyle: {
      w: "200px",
      bg: "primary",
      borderRadius: 10,
      p: 4,
      color: "buttonText",
    },
  },
  Group: {
    baseStyle: {
      w: "300px",
      bg: "primary",
      borderRadius: 10,
      p: 4,
      color: "buttonText",
    },
  },
};