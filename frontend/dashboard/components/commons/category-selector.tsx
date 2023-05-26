import { useTwitchGameCategories } from '@/services/api';
import { Autocomplete, Loader, useMantineTheme } from '@mantine/core';
import React, { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface Props {
  channelId: string;
}

const CategorySelector = ({ channelId }: Props) => {
  const timeoutRef = useRef<number>(-1);
  const [category, setCategory] = useState('');
  const [loading, setLoading] = useState(false);

  const { t } = useTranslation('twitch');
  const theme = useMantineTheme();
  const categories = useTwitchGameCategories(category, channelId);

  const handleChange = (val: string) => {
    window.clearTimeout(timeoutRef.current);
    setCategory(val);

    if (val.trim().length === 0) {
      setLoading(false);
    } else {
      setLoading(true);
      timeoutRef.current = window.setTimeout(() => {
        setLoading(false);
        /*TODO: mutation for categories */
      }, 1000);
    }
  };
  return (
    <div>
      <div>
        <input placeholder="Search for a category" />
      </div>
      <div>
        <div></div>
      </div>
    </div>
  );
};

export default CategorySelector;
