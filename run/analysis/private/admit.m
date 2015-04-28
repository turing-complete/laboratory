function good = admit(indices)
  [nn, ni] = size(indices);

  good = false(nn, 1);
  indices = bitshift(bitshift(indices, 32), -32);
  indices = int32(indices);

  for i = 1:nn
    I = repmat(indices(i, :), ni, 1) - int32(eye(ni));
    I(any(I < 0, 1), :) = [];
    good(i) = all(ismember(I, indices, 'rows'));
  end
end
