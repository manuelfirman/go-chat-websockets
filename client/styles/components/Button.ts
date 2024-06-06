import { extendTheme } from '@chakra-ui/react';
import { colors } from '../../styles/colors';
import { styles } from '../../styles';
import { components } from '../../styles/components';


const theme = extendTheme({colors, styles, components});

export default theme;