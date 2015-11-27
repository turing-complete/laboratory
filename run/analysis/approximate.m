function approximate(grid)
  use('Interaction');

  filename = locate('approximate');
  solution = h5read(filename, '/solution');
  surrogate = solution.Surrogate;

  ni = double(surrogate.Inputs);
  no = double(surrogate.Outputs);
  indices = reshape(surrogate.Indices{1}, ni, [])';
  surpluses = reshape(surrogate.Surpluses{1}, no, [])';

  Plot.figure(1200, 400);
  for i = 1:no
    subplot(no, 1, i);
    semilogy(abs(surpluses(:, i)));
    Plot.title('Output %d', i-1);
  end

  nn = size(indices, 1);

  nu = size(unique(indices, 'rows'), 1);
  if nu ~= nn
    warning('found %d nonunique indices out of %d', nn-nu, nn);
  end

  if nargin == 0, return; end

  nodes = zeros(nn, ni);
  nodes(:) = feval(grid, indices);

  switch ni
  case 2
    Plot.figure(600, 600);
    Plot.line(nodes(:, 1), nodes(:, 2), 'discrete', true, 'style', ...
      {'MarkerFaceColor', 'k', 'MarkerEdgeColor', 'None', 'MarkerSize', 1});
    Plot.title('Nodes');
  end
end

function nodes = open(indices)
  nodes = zeros(size(indices));
  for i = 1:numel(indices)
    nodes(i) = double(bitshift(indices(i), -32)+1) / ...
      double(bitshift(uint64(2), bitshift(bitshift(indices(i), 32), -32)));
  end
end

function nodes = closed(indices)
  nodes = zeros(size(indices));
  for i = 1:numel(indices)
    level = bitshift(bitshift(indices(i), 32), -32);
    if level == 0
      nodes(i) = 0.5;
    else
      nodes(i) = double(bitshift(indices(i), -32)) / ...
        double(bitshift(uint64(2), level-1));
    end
  end
end
